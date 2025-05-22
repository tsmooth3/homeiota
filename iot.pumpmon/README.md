# iot.pumpmon

Pump Monitor for Home IoT System

This component is a Python-based monitor for well pump devices, designed to run on microcontrollers compatible with CircuitPython. It measures current draw, detects pump activity, and reports status to a central API and Slack.

## Features
- Reads analog current sensor to detect pump on/off and running state
- Sends runtime and current data to a Go-based Home API
- Sends alerts/notifications to Slack
- Periodic heartbeat reporting
- Visual feedback via NeoPixel

## Requirements
- CircuitPython-compatible board (e.g., Adafruit Feather)
- Analog current sensor connected to A1
- NeoPixel onboard or attached
- WiFi credentials and API/Slack keys set as environment variables or in `settings.toml`
- Required CircuitPython libraries in `lib/` (see below)

## Setup
1. Copy `pumpmon.py`, `settings.toml`, and the `lib/` folder to your device.
2. Edit `settings.toml` to set your WiFi, Slack, and API credentials.
3. Ensure the following libraries are present in `lib/`:
   - `adafruit_requests.mpy`
   - `neopixel.mpy`
   - `adafruit_connection_manager.mpy`
   - `adafruit_ht16k33/` (if needed)
4. Connect your current sensor to analog pin A1.
5. Power on the device.

## Environment Variables / Settings
- `slack_key`: Slack webhook URL
- `gohomeapi_key`: Base URL for the Go Home API
- `wifi_ssid`: WiFi SSID
- `wifi_pass`: WiFi password

These can be set in `settings.toml` or as environment variables.

## How It Works
- On startup, connects to WiFi and sends a Slack notification and heartbeat.
- Continuously samples current sensor to detect pump state.
- Sends data to the API and Slack as appropriate.
- Blinks NeoPixel for status and error feedback.

## Main Files
- `pumpmon.py`: Main monitoring script
- `settings.toml`: Configuration file
- `lib/`: Required CircuitPython libraries

## Example Usage
No manual interaction is required. The script runs automatically on boot.

## License
See [../LICENSE](../LICENSE) for details.
