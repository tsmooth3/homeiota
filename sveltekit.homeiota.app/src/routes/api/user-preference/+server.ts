import { json } from '@sveltejs/kit';
import { db } from '$lib/server/database';

// query to get user preferences for current user session
export async function GET({ cookies }) {
  const sessionId = cookies.get('session');
  if (!sessionId) {
    return json({ error: 'Session not found' }, { status: 401 });
  }

  const session = await db.session.findUnique({
    where: { id: sessionId },
    include: { user: true }
  });

  if (!session || session.expiresAt < new Date()) {
    return json({ error: 'Session expired' }, { status: 401 });
  }

  // Fetch all alert preferences for this user
  const alertPreferences = await db.alertPreference.findMany({
    where: { userId: session.user.id }
  });

  return json({
    alertPreferences
  });
}
