## Weather Api

Microservice which can return today, hourly and daily weather forecast by http apies.
This service has 4 apies:
1. Get today weather forecast
2. Get hourly weather forecast (3 days)
3. Get daily weather forecast (10 days)
4. Get place id by your city details

For details about apies, see src/main.py You can find controllers for all apies with documentation.

### Requirements
- Docker Engine

### Steps to run

1. Build Dockerfile image
```bash
    docker build -t weather-api:v1.0.0 .
```
2. Start image
```bash
    docker run -p 8000:8000 {image_id}
```