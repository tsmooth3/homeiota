// See https://kit.svelte.dev/docs/types#app
// for information about these interfaces
declare global {
  namespace App {
    // interface Error {}
    interface Locals {
      sessionId?: string;
      user?: {
        id: string;
        email: string;
        name: string;
        alertPreferences?: {
          id: string;
          thresholds: {
            temperature: {
              enabled: boolean;
              value: number;
            };
            offline: {
              enabled: boolean;
              gracePeriod: number;
            };
          };
        };
      };
    }
    // interface PageData {}
    // interface Platform {}
  }
}

export {}; 