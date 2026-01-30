FROM python:3.9-slim

WORKDIR /app

COPY alert-service/requirements.txt .
RUN pip install -r requirements.txt

COPY alert-service/ .

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
