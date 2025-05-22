# ESP32-C3/S3 CircuitPython Temperature Monitor

This project is a CircuitPython script for the ESP32-C3 or ESP32-S3 microcontroller. It supports:
- SHT4x temperature/humidity sensor (I2C)
- Optional 4-digit 7-segment display (HT16K33, I2C)
- Optional NeoPixel (on-board or external)

The script reads temperature, displays it (if display is present), sends readings to a GoHome API, and sends notifications to Slack. It also blinks the NeoPixel (if present) to indicate network activity.

## Features
- Reads temperature from SHT4x sensor
- Displays temperature on 7-segment display (if detected)
- Connects to WiFi and sends data to GoHome API
- Sends startup and error notifications to Slack
- Blinks NeoPixel for Slack (orange) and GoHomeAPI (teal) events
- Prints IP address and WiFi signal strength on startup
- All hardware is optional except the SHT4x sensor (required for temperature readings)

## Hardware Requirements
- ESP32-C3 or ESP32-S3 board running CircuitPython
- SHT4x temperature/humidity sensor (I2C)
- HT16K33 4-digit 7-segment display (I2C, optional)
- NeoPixel (on-board or external, optional)

## Setup
1. **Install CircuitPython** on your ESP32 board.
2. **Copy the following files to your CIRCUITPY drive:**
   - `code.py` (this script)
   - All required libraries in the `lib/` folder:
     - `adafruit_ht16k33/`
     - `adafruit_sht4x.mpy`
     - `adafruit_requests.mpy`
     - `adafruit_connection_manager.mpy`
     - `neopixel.mpy`
3. **Configure WiFi and API settings:**
   - Create a `settings.toml` file or set environment variables for:
     - `wifi_ssid` and `wifi_pass` (WiFi credentials)
     - `gohomeapi_key` (base URL for GoHome API)
     - `slack_key` (Slack webhook URL, optional)
     - `TEMP_MON_LOCATION` (location label, optional)

## Usage
- On boot, the device connects to WiFi, prints its IP and signal strength, reads the temperature, displays it (if display is present), and sends it to the GoHome API.
- The temperature is updated and sent every 30 seconds (configurable via `MEASURE_INTERVAL`).
- If configured, a Slack message is sent on startup with temperature, IP, and signal strength.
- NeoPixel blinks orange for Slack, teal for GoHomeAPI events.

## Example `settings.toml`
```toml
wifi_ssid = "YourWiFiSSID"
wifi_pass = "YourWiFiPassword"
gohomeapi_key = "http://your.api.endpoint"
slack_key = "https://hooks.slack.com/services/your/webhook/url"
TEMP_MON_LOCATION = "freezer"
```

## CircuitPython Libraries
Download the required libraries from the [Adafruit CircuitPython Bundle](https://circuitpython.org/libraries).

## Pin Connections
- **SHT4x**: Connect to I2C (SCL/SDA)
- **HT16K33 7-segment**: Connect to I2C (SCL/SDA)
- **NeoPixel**: On-board or connect to a supported pin

## License
MIT
