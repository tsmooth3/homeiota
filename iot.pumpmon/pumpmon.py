import os
import time
import adafruit_requests as requests
import analogio
import board
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
GOHOMEAPI = f"{GOHOMEAPIBASE}/pumpmon"
GOHOMEHEARTBEAT = f"{GOHOMEAPIBASE}/heartbeat"

analog_pin = analogio.AnalogIn(board.A1)
pixel = neopixel.NeoPixel(board.NEOPIXEL, 1)
pixel.brightness = 0.3
pumpOnAmps = 1.7
pumpRunAmps = 6.0


def get_ip():
    # WiFi setup
    wifi.radio.enabled = False
    time.sleep(1)
    wifi.radio.enabled = True
    time.sleep(1)
    wifi.radio.connect(WIFISSID, WIFIPASS)
    pool = socketpool.SocketPool(wifi.radio)
    https = requests.Session(pool, ssl.create_default_context())
    return wifi, https

wifi, https = get_ip()

def blink(duration, color=(0, 0, 0)):
    for i in range(duration):
        pixel.fill(color)
        time.sleep(0.1)
        pixel.fill((0, 0, 0))
        time.sleep(0.1)

def send_slack_message(wifi, https, message):
    try:
        if wifi.radio.connected == False:
            wifi, https = get_ip()
        response = https.post(SLACK_URL, json={"text": message})
        print(f"Slack response code: {response.status_code}")
        blink(1, color=(0, 255, 0))
    except Exception as e:
        print(f"Error sending Slack message: {e}")
        blink(5, color=(255, 0, 0))

def send_gohomeapi_data(wifi, https, runtime, current, isLow):
    try:
        if wifi.radio.connected == False:
            wifi, https = get_ip()
        payload = {
            "run_time": runtime,
            "current": current,
            "low_current": isLow
        }
        response = https.post(GOHOMEAPI, json=payload)
        print(f"GoHomeAPI response code: {response.status_code}")
        blink(1, color=(0, 255, 0))
    except Exception as e:
        print(f"Error sending GoHomeAPI data: {e}")
        blink(5, color=(255, 0, 0))

def send_heartbeat(wifi, https):
    try:
        if wifi.radio.connected == False:
            wifi, https = get_ip()
        payload = {
            "device_id": "pump"
        }
        response = https.post(GOHOMEHEARTBEAT, json=payload)
        print(f"Heartbeat response code: {response.status_code}")
        blink(2, color=(255, 0, 255))
    except Exception as e:
        print(f"Error sending Heartbeat data: {e}")
        blink(5, color=(255, 0, 0))

def readCurrent():
    try:
        num_samples = 30 # Number of samples to take
        adc_value = 0

        for i in range(num_samples):
            adc_value += analog_pin.value
            time.sleep(0.01)

        # Calculate the average value
        adc_value /= num_samples 
        # Convert ADC value to voltage
        # 8191 is the maximum value for 13-bit ADC, 3.3V is the reference voltage
        voltage = (adc_value / 8191.0) * 4.94 # 4.94V is the reference voltage for the current sensor
        # Convert voltage to current (in Amps)  
         # 0.05V/A is the sensitivity of the current sensor
        current = voltage / 0.05
        return current
    except Exception as e:
        print(f"Error reading temperature: {e}")
        return None

send_slack_message(wifi, https, f"Starting up well pump monitor")
send_heartbeat(wifi, https)

startuptime = time.monotonic()
currentTime = time.monotonic()
heartbeatTime = time.monotonic()
pumpOn = False
while True:
    current = readCurrent()
    
    if heartbeatTime + 60 < time.monotonic():
        send_heartbeat(wifi, https)
        heartbeatTime = time.monotonic()
    else:
        blink(2, color=(0, 255, 255))

    if current > pumpOnAmps:
        if not pumpOn:
            pumpOn = True
            startuptime = time.monotonic()
            print(f"Pump On: {current:.2f}")
        currentTime = time.monotonic()
        runtime =  int(currentTime - startuptime)        
        print(f"Pump On: {runtime} : {current:.2f}")
        if current > pumpRunAmps:
            send_gohomeapi_data(wifi, https, runtime, current, False)
        else:
            send_gohomeapi_data(wifi, https, runtime, current, True)
    else:
        print(f"Pump Off: {current:.2f}")
        if pumpOn:
            currentTime = time.monotonic()
            runtime =  int(currentTime - startuptime)
            print(f"Pump Off: {runtime} : {current:.2f}")
            pumpOn = False
            send_gohomeapi_data(wifi, https, runtime, current, False)

    time.sleep(2)