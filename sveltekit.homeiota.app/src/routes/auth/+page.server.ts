import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { db } from '$lib/server/database';
import bcrypt from 'bcryptjs';
import { createSession } from '$lib/server/session';

export const actions: Actions = {
  default: async ({ request, cookies }) => {
    const data = await request.formData();
    const email = data.get('email') as string;
    const password = data.get('password') as string;
    const name = data.get('name') as string;
    const isLogin = data.get('isLogin') === 'true';

    console.log('Auth request:', { email, isLogin, hasName: !!name });

    if (!email || !password) {
      return fail(400, { error: 'Email and password are required' });
    }

    try {
      if (isLogin) {
        // Login
        const user = await db.user.findUnique({
          where: { email },
          include: {
            alertPreferences: true
          }
        });
        if (!user) {
          return fail(400, { error: 'Invalid email or password' });
        }

        const validPassword = await bcrypt.compare(password, user.password);
        if (!validPassword) {
          return fail(400, { error: 'Invalid email or password' });
        }

        // Create session
        const session = await createSession(user.id);
        console.log('Login - Session created:', session.id);
        
        // Set cookie with more permissive settings for development
        cookies.set('session', session.id, {
          path: '/',
          httpOnly: true,
          sameSite: 'lax',
          secure: false,
          maxAge: 60 * 60 * 24 * 7 // 1 week
        });
        console.log('Login - Session cookie set');

        return redirect(303, '/');
      }

      // Registration
      if (!name) {
        return fail(400, { error: 'Name is required' });
      }

      const existingUser = await db.user.findUnique({ where: { email } });
      if (existingUser) {
        return fail(400, { error: 'Email already registered' });
      }

      const hashedPassword = await bcrypt.hash(password, 10);

      // First create the user
      const user = await db.user.create({
        data: {
          email,
          name,
          password: hashedPassword
        }
      });

      

      // Create session
      const session = await createSession(user.id);
      console.log('Registration - Session created:', session.id);
      
      // Set cookie with more permissive settings for development
      cookies.set('session', session.id, {
        path: '/',
        httpOnly: true,
        sameSite: 'lax',
        secure: false,
        maxAge: 60 * 60 * 24 * 7 // 1 week
      });
      console.log('Registration - Session cookie set');

      return redirect(303, '/');
    } catch (error) {
      console.error('Auth error:', error);
      return fail(500, { error: 'An error occurred during authentication' });
    }
  }
}; 