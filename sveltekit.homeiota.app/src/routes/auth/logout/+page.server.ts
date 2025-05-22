import { redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { db } from '$lib/server/database';

export const actions: Actions = {
  default: async ({ cookies }) => {
    const sessionId = cookies.get('session');
    if (sessionId) {
      await db.session.delete({
        where: { id: sessionId }
      });
      cookies.delete('session', { path: '/' });
    }
    throw redirect(303, '/');
  }
}; 