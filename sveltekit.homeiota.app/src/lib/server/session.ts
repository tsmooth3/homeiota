import { db } from './database';
import type { Handle } from '@sveltejs/kit';
import { Prisma } from '@prisma/client';

// Extend App.Locals for type safety
// (You may want to move this to a global types file)
declare module '@sveltejs/kit' {
  interface Locals {
    sessionId?: string;
    user?: any; // Replace 'any' with your User type if available
  }
}

export async function getSession(locals: App.Locals) {
  const sessionId = locals.sessionId;
  if (!sessionId) return null;

  const session = await db.session.findUnique({
    where: { id: sessionId },
    include: { user: true }
  });

  if (!session || session.expiresAt < new Date()) {
    return null;
  }

  return session;
}

export async function createSession(userId: string) {
  const expiresAt = new Date();
  expiresAt.setDate(expiresAt.getDate() + 7); // 1 week

  return db.session.create({
    data: {
      userId,
      expiresAt
    }
  });
}

export const handle: Handle = async ({ event, resolve }) => {
  const sessionId = event.cookies.get('session');
  
  if (sessionId) {
    const session = await db.session.findUnique({
      where: { id: sessionId },
      include: { user: true }
    });

    if (session && session.expiresAt > new Date()) {
      event.locals.sessionId = sessionId;
      event.locals.user = session.user;
    } else {
      // Delete expired session, ignore if it doesn't exist
      try {
        await db.session.delete({
          where: { id: sessionId }
        });
      } catch (error) {
        if (
          error instanceof Prisma.PrismaClientKnownRequestError &&
          error.code === 'P2025'
        ) {
          // Session already deleted, ignore
        } else {
          throw error;
        }
      }
      event.cookies.delete('session', { path: '/' });
    }
  }

  const response = await resolve(event);
  return response;
}; 