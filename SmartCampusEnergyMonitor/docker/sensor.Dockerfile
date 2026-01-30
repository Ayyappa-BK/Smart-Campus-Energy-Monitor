FROM python:3.9-slim

WORKDIR /app

COPY sensor-simulator/requirements.txt .
# I'll create it or just install inline
RUN pip install grpcio grpcio-tools protobuf

COPY sensor-simulator/ .
# Copy proto file to generate valid imports if needed, or rely on pre-generated
# In our flow we generated code into sensor-simulator/ so we just copy it.

CMD ["python", "simulator.py"]
