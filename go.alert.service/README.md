# go.alert.service

Go Alert Service for Home IoT System

This service monitors device data (such as pump run times and temperatures) and sends alerts/notifications to users via Gotify. It queries PostgreSQL databases for recent device activity and user alert preferences, and triggers notifications when thresholds are exceeded or devices go offline.

## Features
- Sends alerts to Gotify based on user preferences
- Monitors pump run times, temperature readings, and device heartbeats
- Supports offline/device-down detection
- Configurable thresholds per user/location
- Connects to multiple PostgreSQL databases

## Requirements
- Go 1.18+
- PostgreSQL databases for device data and user preferences
- Gotify server for notifications
- Environment variables for configuration

## Environment Variables
- `GOTIFY_URL`: Base URL for Gotify server (e.g., `http://gotify.example.com`)
- `HOMEIOTA_URL`: URL for the Home IoT dashboard (used in alert messages)
- `GOHOME_DB_URL`: Connection string for the Go Home API database
- `HOMEIOTA_DB_URL`: Connection string for the Home IoT user/alert preferences database

## Setup & Usage

### Build and Run with Docker (Recommended)
```bash
cd go.alert.service
docker build -t homeiota-alert-service .
docker run --env-file .env homeiota-alert-service
```

### Or Run Locally
```bash
cd go.alert.service
go run main.go
```

## How It Works
- On execution, connects to the configured databases
- Fetches user alert preferences and recent device data
- Checks for threshold violations and device offline status
- Sends Gotify notifications for:
  - Pump current anomalies
  - High temperature readings
  - Device offline/heartbeat missing

## Main Files
- `main.go`: Main application logic
- `Dockerfile`: Containerization support
- `go.mod`, `go.sum`: Go module dependencies

## License
See [../LICENSE](../LICENSE) for details.
