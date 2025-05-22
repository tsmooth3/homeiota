import { fail, redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';
import { db } from '$lib/server/database';

export const load: PageServerLoad = async ({ cookies }) => {
  const sessionId = cookies.get('session');
  if (!sessionId) {
    throw redirect(303, '/auth');
  }

  const session = await db.session.findUnique({
    where: { id: sessionId },
    include: { user: true }
  });

  if (!session || session.expiresAt < new Date()) {
    throw redirect(303, '/auth');
  }

  // Fetch all alert preferences for this user
  const alertPreferences = await db.alertPreference.findMany({
    where: { userId: session.user.id }
  });

  return {
    user: session.user,
    alertPreferences
  };
};

export const actions: Actions = {
  default: async ({ request, cookies }) => {
    const sessionId = cookies.get('session');
    if (!sessionId) {
      throw redirect(303, '/auth');
    }

    const session = await db.session.findUnique({
      where: { id: sessionId },
      include: { user: true }
    });

    if (!session || session.expiresAt < new Date()) {
      throw redirect(303, '/auth');
    }

    let name = '';
    let email = '';
    let gotifyToken = '';
    let uiAlertPreferences: { name: string; threshold: number; enabled: boolean; offlineThreshold?: number }[] = [];

    // Always use formData
    const data = await request.formData();
    name = data.get('name') as string;
    email = data.get('email') as string;
    gotifyToken = data.get('gotifyToken') as string;
    const sensorsJson = data.get('uiAlertPreferences');
    uiAlertPreferences = sensorsJson ? JSON.parse(sensorsJson as string) : [];

    try {
      // Update user profile if present
      if (name && email) {
        await db.user.update({
          where: { id: session.user.id },
          data: {
            name,
            email,
            gotifyToken
          }
        });
      }

      // Upsert alert preferences if present
      if (uiAlertPreferences && Array.isArray(uiAlertPreferences)) {
        for (const sensor of uiAlertPreferences) {
          await db.alertPreference.upsert({
            where: {
              userId_location: {
                userId: session.user.id,
                location: sensor.name
              }
            },
            update: {
              threshold: sensor.threshold,
              enabled: sensor.enabled,
              offlineThreshold: sensor.offlineThreshold ?? null
            },
            create: {
              userId: session.user.id,
              location: sensor.name,
              threshold: sensor.threshold,
              enabled: sensor.enabled,
              offlineThreshold: sensor.offlineThreshold ?? null
            }
          });
        }
      }

      // TODO: Handle offline alert preferences if you want to store them in a table

      // Ensure the session cookie is still set
      cookies.set('session', sessionId, {
        path: '/',
        httpOnly: true,
        sameSite: 'strict',
        secure: typeof process !== 'undefined' && process.env.NODE_ENV === 'production',
        maxAge: 60 * 60 * 24 * 7 // 1 week
      });

      return { success: true };
    } catch (error) {
      console.error('Failed to update settings:', error);
      return fail(500, { error: 'Failed to update settings' });
    }
  }
}; 