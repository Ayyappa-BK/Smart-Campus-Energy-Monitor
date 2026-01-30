from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import boto3
import os
import time
import json
from datetime import datetime

app = FastAPI()

# S3 Configuration
S3_BUCKET = os.environ.get("S3_BUCKET_NAME", "smart-campus-alerts-dev")
AWS_REGION = os.environ.get("AWS_REGION", "us-east-1")

# Initialize S3 Client (assuming credentials are in environment or ~/.aws/credentials)
# For local dev without creds, this might fail unless mocked.
# We will wrap in try-except for robustness during demo.
s3_client = boto3.client('s3', region_name=AWS_REGION)

class AlertPayload(BaseModel):
    building_id: str
    floor_id: str
    wattage: float
    timestamp: int
    message: str

@app.post("/alert")
async def receive_alert(alert: AlertPayload):
    print(f"CRITICAL ALERT RECEIVED: {alert}")
    
    # Create a log entry
    log_entry = alert.dict()
    log_entry["received_at"] = datetime.utcnow().isoformat()
    
    file_name = f"alert_{alert.building_id}_{alert.timestamp}.json"
    
    try:
        s3_client.put_object(
            Bucket=S3_BUCKET,
            Key=f"alerts/{datetime.utcnow().date()}/{file_name}",
            Body=json.dumps(log_entry),
            ContentType='application/json'
        )
        print(f"Alert uploaded to S3: {S3_BUCKET}/{file_name}")
        return {"status": "logged", "s3_key": file_name}
    except Exception as e:
        print(f"Failed to upload to S3: {e}")
        # In a real app we might retry or store locally
        return {"status": "error", "reason": str(e)}

@app.get("/health")
def health():
    return {"status": "ok"}
