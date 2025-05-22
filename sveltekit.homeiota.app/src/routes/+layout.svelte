<script lang="ts">
  import { goto } from '$app/navigation';
  import '../app.css';
  export let data: { isAuthenticated: boolean; user: { name: string } | null };

  function goHome() {
    goto('/');
  }

  function handleKeydown(event) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      goHome();
    }
  }
</script>

<div class="min-h-screen bg-gray-900 text-white">
  <nav class="bg-gray-800 border-b border-gray-700 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <div class="flex items-center">
          <a href="/" class="flex items-center">
            <img class="h-8 w-8" src="/favicon.png" alt="Home IoTa" />
            <span class="ml-2 text-xl font-bold">Home IoTa</span>
          </a>
        </div>
        <div class="flex items-center space-x-4">
          {#if data.isAuthenticated}
            <span class="text-gray-300">Hi, {data.user?.name}</span>
            <a href="/settings" class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Settings</a>
            <form action="/auth/logout" method="POST" class="inline">
              <button type="submit" class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium">
                Logout
              </button>
            </form>
          {:else}
            <a href="/auth" class="text-gray-300 hover:text-white px-3 py-2 rounded-md text-sm font-medium">Login</a>
          {/if}
        </div>
      </div>
    </div>
  </nav>

  <main class="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
    <slot />
  </main>
</div> 