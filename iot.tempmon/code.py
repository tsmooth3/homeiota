import os
import time
import board
import adafruit_ht16k33.segments
import adafruit_requests as requests
import adafruit_sht4x
import wifi
import socketpool
import ssl
import neopixel
wifi.radio.tx_power = 13

# Constants
SLACK_URL = os.getenv("slack_key")
GOHOMEAPIBASE = os.getenv("gohomeapi_key")
WIFISSID = os.getenv("wifi_ssid")
WIFIPASS = os.getenv("wifi_pass")
GOHOMEAPI = f"{GOHOMEAPIBASE}/tempmon"
TEMP_MON_LOCATION = os.getenv("TEMP_MON_LOCATION", "tempmon")
MEASURE_INTERVAL = 30  # seconds

sht = adafruit_sht4x.SHT4x(board.STEMMA_I2C())
sht.mode = adafruit_sht4x.Mode.NOHEAT_HIGHPRECISION

# Try to initialize the matrix display
try:
    i2c = board.I2C()
    matrix = adafruit_ht16k33.segments.Seg7x4(i2c)
    matrix.brightness = 0.3
    matrix.print("ello")
    matrix.show()
except Exception as e:
    print(f"Matrix display not detected: {e}")
    matrix = None

# Try to initialize the NeoPixel (on pin board.NEOPIXEL, 1 pixel)
try:
    pixel = neopixel.NeoPixel(board.NEOPIXEL, 1)
    pixel.brightness = 0.1
except Exception as e:
    print(f"NeoPixel not detected: {e}")
    pixel = None

# WiFi setup
wifi.radio.connect(WIFISSID, WIFIPASS)
time.sleep(1)
# print the ip address
print(f"Connected to WiFi. IP address: {wifi.radio.ipv4_address}")
# print the signal strength
print(f"Signal strength: {wifi.radio.ap_info.rssi} dBm")
# create a socket pool and https session
pool = socketpool.SocketPool(wifi.radio)
https = requests.Session(pool, ssl.create_default_context())

def blink_pixel(color=(0, 255, 0), duration=0.1, times=2):
    if pixel:
        for _ in range(times):
            pixel[0] = color
            time.sleep(duration)
            pixel[0] = (0, 0, 0)
            time.sleep(duration)

def send_slack_message(message):
    try:
        response = https.post(SLACK_URL, json={"text": message})
        print(f"Slack response code: {response.status_code}")
        if pixel:
            #blink orange for Slack
            blink_pixel((255, 165, 0)) 
    except Exception as e:
        print(f"Error sending Slack message: {e}")

def send_gohomeapi_data(temp):
    try:
        payload = {
            "value": temp,
            "location": TEMP_MON_LOCATION
        }
        response = https.post(GOHOMEAPI, json=payload)
        print(f"GoHomeAPI response code: {response.status_code}")
        if pixel:
            #blink teal for GoHomeAPI 
            blink_pixel((0, 255, 255)) 
    except Exception as e:
        print(f"Error sending GoHomeAPI data: {e}")

def readTemp():
    try:
        temp, rh = sht.measurements
        tempF = (temp * 1.8) + 32
        return tempF
    except Exception as e:
        print(f"Error reading temperature: {e}")
        return None

# send a startup message to Slack - include the current temperature, ip address, and signal strength
send_slack_message(f"Starting up {TEMP_MON_LOCATION} temperature monitor: {readTemp():.2f} F ip: {wifi.radio.ipv4_address} signal: {wifi.radio.ap_info.rssi} dBm")
while True:
    tempF = readTemp()
    send_gohomeapi_data(tempF)
    print(f"temp: {tempF:.2f} F")
    if matrix:
        matrix.fill(0)
        matrix.print(f"{tempF:05.1f}")
        matrix.show()
    time.sleep(MEASURE_INTERVAL)
