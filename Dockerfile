FROM golang:alpine as builder

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

RUN go build -o /app/weather-api

FROM builder AS run-test-stage
RUN go test -v ./...

FROM alpine

WORKDIR /app

COPY --from=builder /app/weather-api /app/weather-api

CMD ["./weather-api"]