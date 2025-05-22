import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
  return {
    isAuthenticated: !!locals.user,
    user: locals.user
  };
}; 