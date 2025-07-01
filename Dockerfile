FROM python:3.13-alpine

WORKDIR /app

COPY . .

# Install dependencies for python instruments
RUN apk add --no-cache \
        build-base \
        musl-dev \
        linux-headers \
        python3-dev

# Download and install python dependencies
RUN pip3 install --no-cache-dir -r requirements.txt

CMD ["uvicorn", "src.main:app", "--host", "0.0.0.0", "--port", "8000", "--log-level", "info"]