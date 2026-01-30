package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"sync"
	"encoding/json"
	"bytes"

	"google.golang.org/grpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"

	pb "aggregator-service/pb"
)

// Metrics
var (
	wattageGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "energy_current_wattage",
			Help: "Current wattage consumption per building/floor",
		},
		[]string{"building_id", "floor_id"},
	)
)

func init() {
	prometheus.MustRegister(wattageGauge)
}

type server struct {
	pb.UnimplementedEnergySensorServer
	mu sync.Mutex
	// Store history for anomaly detection (simple moving average)
	// Map: building_id -> []wattage
	history map[string][]float64
}

func (s *server) StreamEnergyData(stream pb.EnergySensor_StreamEnergyDataServer) error {
	log.Println("New stream connection established")
	for {
		reading, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{
				Success: true,
				Message: "Stream closed",
			})
		}
		if err != nil {
			log.Printf("Error receiving data: %v", err)
			return err
		}

		// Process data
		s.processReading(reading)
	}
}

func (s *server) processReading(r *pb.EnergyReading) {
	log.Printf("Received: Building=%s Floor=%s Wattage=%.2f Time=%d", 
		r.BuildingId, r.FloorId, r.CurrentWattage, r.Timestamp)

	// Update Prometheus
	wattageGauge.WithLabelValues(r.BuildingId, r.FloorId).Set(r.CurrentWattage)

	// Anomaly Detection (Simple Moving Average)
	s.detectAnomaly(r)
}

func (s *server) detectAnomaly(r *pb.EnergyReading) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := r.BuildingId
	windowSize := 10 // Last 10 readings
	
	if s.history == nil {
		s.history = make(map[string][]float64)
	}

	readings := s.history[key]
	readings = append(readings, r.CurrentWattage)
	if len(readings) > windowSize {
		readings = readings[1:]
	}
	s.history[key] = readings

	if len(readings) < windowSize {
		return // Not enough data
	}

	// Calculate Avg and StdDev
	var sum float64
	for _, v := range readings {
		sum += v
	}
	mean := sum / float64(len(readings))

	var variance float64
	for _, v := range readings {
		variance += math.Pow(v-mean, 2)
	}
	stdDev := math.Sqrt(variance / float64(len(readings)))

	// Z-Score
	if stdDev > 0 {
		zScore := (r.CurrentWattage - mean) / stdDev
		if math.Abs(zScore) > 3 { // Threshold > 3 sigma
			log.Printf("ANOMALY DETECTED! Building: %s, Value: %.2f, Z-Score: %.2f", key, r.CurrentWattage, zScore)
			go s.sendAlert(r, zScore)
		}
	}
}

func (s *server) sendAlert(r *pb.EnergyReading, zScore float64) {
	alertPayload := map[string]interface{}{
		"building_id": r.BuildingId,
		"floor_id":    r.FloorId,
		"wattage":     r.CurrentWattage,
		"timestamp":   r.Timestamp,
		"message":     fmt.Sprintf("Z-Score Anomaly: %.2f", zScore),
	}
	
	jsonBody, _ := json.Marshal(alertPayload)
	alertServiceURL := "http://alert-service:8000/alert" 

	resp, err := http.Post(alertServiceURL, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Failed to send alert: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Alert sent, status: %s", resp.Status)
}


func main() {
	// Start Prometheus Metrics Server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Metrics server listening on :2112")
		log.Fatal(http.ListenAndServe(":2112", nil))
	}()

	// Start get gRPC Server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEnergySensorServer(s, &server{})
	
	log.Println("Aggregator Service listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
