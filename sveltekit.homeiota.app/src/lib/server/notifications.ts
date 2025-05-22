import got from 'got';

const GOTIFY_URL = process.env.GOTIFY_URL;
const GOTIFY_TOKEN = process.env.GOTIFY_TOKEN;

export async function sendNotification(userId: string, message: string, priority: number = 5) {
  try {
    await got.post(`${GOTIFY_URL}/message`, {
      json: {
        message,
        priority,
        extras: {
          'client::display': {
            contentType: 'text/markdown'
          }
        }
      },
      headers: {
        'X-Gotify-Key': GOTIFY_TOKEN
      }
    });
  } catch (error) {
    console.error('Failed to send notification:', error);
    throw error;
  }
}

// Helper function to send temperature alerts
export async function sendTemperatureAlert(deviceName: string, temperature: number, threshold: number) {
  const message = `🌡️ Temperature Alert\n\nDevice: ${deviceName}\nCurrent: ${temperature}°F\nThreshold: ${threshold}°F`;
  await sendNotification('temperature', message, 7); // Higher priority for temperature alerts
}

// Helper function to send offline alerts
export async function sendOfflineAlert(deviceName: string, minutesOffline: number) {
  const message = `⚠️ Device Offline\n\nDevice: ${deviceName}\nOffline for: ${Math.round(minutesOffline)} minutes`;
  await sendNotification('offline', message, 8); // Highest priority for offline alerts
} 