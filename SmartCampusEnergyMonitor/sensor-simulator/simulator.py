import grpc
import time
import random
import os
import sys

# Add current directory to path to find generated modules
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

import sensor_pb2
import sensor_pb2_grpc

def generate_readings():
    """Generates random energy readings with occasional spikes."""
    buildings = ["Engineering", "Science", "Library", "DormA"]
    
    while True:
        building = random.choice(buildings)
        floor = f"Floor-{random.randint(1, 4)}"
        
        # Normal wattage 500-1000, Spike 5000+
        is_spike = random.random() < 0.05 # 5% chance of spike
        
        if is_spike:
            wattage = random.uniform(4000, 6000)
            print(f"Generating SPIKE for {building}")
        else:
            wattage = random.uniform(500, 1000)
            
        timestamp = int(time.time())
        
        yield sensor_pb2.EnergyReading(
            building_id=building,
            floor_id=floor,
            current_wattage=wattage,
            voltage=120.0, # mostly constant
            timestamp=timestamp
        )
        
        time.sleep(0.5)

def run():
    print("Starting Sensor Simulator...")
    # Wait for Aggregator to be ready (simulate basic retry)
    time.sleep(5) 
    
    target = os.environ.get("AGGREGATOR_HOST", "localhost:50051")
    with grpc.insecure_channel(target) as channel:
        stub = sensor_pb2_grpc.EnergySensorStub(channel)
        
        try:
            responses = stub.StreamEnergyData(generate_readings())
            print("Stream started.")
            # The server closes the stream with a StreamResponse
            # But since this is client-streaming/bidi implied by loop, we wait?
            # Actually the rpc is "stream EnergyReading) returns (StreamResponse)"
            # So this call returns a single response when we are done.
            print(f"Server response: {responses}")
            
        except grpc.RpcError as e:
            print(f"gRPC Error: {e}")

if __name__ == '__main__':
    run()
