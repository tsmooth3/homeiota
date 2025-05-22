export interface User {
  id: string;
  email: string;
  name: string;
  phone?: string;
  alertPreferences: AlertPreferences;
  createdAt: Date;
  updatedAt: Date;
}

export interface AlertPreferences {
  email: boolean;
  sms: boolean;
  push: boolean;
  thresholds: {
    temperature: {
      enabled: boolean;
      value: number;
    };
    offline: {
      enabled: boolean;
      gracePeriod: number; // minutes
    };
  };
}

export interface Alert {
  id: string;
  userId: string;
  type: 'temperature' | 'offline';
  deviceId: string;
  message: string;
  status: 'pending' | 'sent' | 'failed';
  createdAt: Date;
  sentAt?: Date;
}

export interface Device {
  id: string;
  name: string;
  type: 'pump' | 'temperature';
  status: 'online' | 'offline';
  lastSeen: Date;
  currentValue?: number;
  threshold?: number;
} 