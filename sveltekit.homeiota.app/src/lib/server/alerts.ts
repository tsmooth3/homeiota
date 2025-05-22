import { db } from './database';
import type { Device } from '@prisma/client';
import { sendTemperatureAlert, sendOfflineAlert } from './notifications';

export async function checkDeviceStatus(device: Device) {
  const users = await db.user.findMany({
    include: {
      alertPreferences: true
    }
  });

  for (const user of users) {
    const prefs = user.alertPreferences;
    
    // Check offline status
    if (prefs.thresholds.offline.enabled) {
      const lastSeen = new Date(device.lastSeen);
      const now = new Date();
      const minutesOffline = (now.getTime() - lastSeen.getTime()) / (1000 * 60);
      
      if (minutesOffline > prefs.thresholds.offline.gracePeriod) {
        await createAlert({
          userId: user.id,
          type: 'offline',
          deviceId: device.id,
          message: `Device ${device.name} has been offline for ${Math.round(minutesOffline)} minutes`
        });
        
        // Send Gotify notification
        await sendOfflineAlert(device.name, minutesOffline);
      }
    }

    // Check temperature threshold
    if (device.type === 'temperature' && 
        prefs.thresholds.temperature.enabled && 
        device.currentValue !== null &&
        device.currentValue > prefs.thresholds.temperature.value) {
      await createAlert({
        userId: user.id,
        type: 'temperature',
        deviceId: device.id,
        message: `Temperature sensor ${device.name} is reading ${device.currentValue}°F (threshold: ${prefs.thresholds.temperature.value}°F)`
      });
      
      // Send Gotify notification
      await sendTemperatureAlert(device.name, device.currentValue, prefs.thresholds.temperature.value);
    }
  }
}

async function createAlert({ userId, type, deviceId, message }: {
  userId: string;
  type: 'temperature' | 'offline';
  deviceId: string;
  message: string;
}) {
  const alert = await db.alert.create({
    data: {
      userId,
      type,
      deviceId,
      message,
      status: 'pending'
    }
  });

  // Send alerts based on user preferences
  const user = await db.user.findUnique({
    where: { id: userId },
    include: { alertPreferences: true }
  });

  if (!user) return;

  const { alertPreferences } = user;
  const promises = [];


  try {
    await Promise.all(promises);
    await db.alert.update({
      where: { id: alert.id },
      data: {
        status: 'sent',
        sentAt: new Date()
      }
    });
  } catch (error) {
    console.error('Failed to send alerts:', error);
    await db.alert.update({
      where: { id: alert.id },
      data: { status: 'failed' }
    });
  }
} 