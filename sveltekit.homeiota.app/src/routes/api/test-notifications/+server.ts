import { json } from '@sveltejs/kit';
import { sendTemperatureAlert, sendOfflineAlert } from '$lib/server/notifications';

export async function POST({ request }) {
  const data = await request.json();
  const { type, deviceId, threshold } = data;

  try {
    if (type === 'temperature') {
      // Use provided deviceId and threshold, use a dummy current value
      await sendTemperatureAlert(
        deviceId || 'Test Device',
        (threshold ?? 80) + 5, // Simulate current value above threshold
        threshold ?? 80
      );
      return json({ success: true, message: 'Temperature alert sent' });
    } else if (type === 'offline') {
      // Test offline alert
      await sendOfflineAlert('Test Device', 10); // 10 minutes offline
      return json({ success: true, message: 'Offline alert sent' });
    } else {
      return json({ error: 'Invalid notification type' }, { status: 400 });
    }
  } catch (error) {
    console.error('Failed to send test notification:', error);
    return json({ error: 'Failed to send notification' }, { status: 500 });
  }
} 