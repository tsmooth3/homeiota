<script lang="ts">
  import { enhance } from '$app/forms';
  import { goto } from '$app/navigation';
  import type { ActionData } from './$types';

  export let form: ActionData;

  let isLogin = true;
  let email = '';
  let password = '';
  let name = '';

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();
    const form = event.target as HTMLFormElement;
    const formData = new FormData(form);
    
    try {
      const response = await fetch(form.action, {
        method: 'POST',
        body: formData,
        redirect: 'follow'
      });
      
      if (response.ok) {
        window.location.href = '/settings';
      } else {
        const result = await response.json();
        if (result.error) {
          form.error = result.error;
        }
      }
    } catch (error) {
      console.error('Form submission error:', error);
    }
  }
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
  <div class="max-w-md w-full space-y-8">
    <div>
      <img class="mx-auto h-16 w-16" src="/favicon.png" alt="Home IoTa" />
      <h2 class="mt-6 text-center text-3xl font-extrabold text-white">
        {isLogin ? 'Sign in to your account' : 'Create a new account'}
      </h2>
    </div>
    <form class="mt-8 space-y-6" method="POST" on:submit={handleSubmit}>
      <input type="hidden" name="isLogin" value={isLogin} />
      <div class="rounded-md shadow-sm -space-y-px">
        {#if !isLogin}
          <div>
            <label for="name" class="sr-only">Name</label>
            <input
              id="name"
              name="name"
              type="text"
              required
              bind:value={name}
              class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-700 bg-gray-800 text-white placeholder-gray-400 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
              placeholder="Name"
            />
          </div>
        {/if}
        <div>
          <label for="email" class="sr-only">Email address</label>
          <input
            id="email"
            name="email"
            type="email"
            required
            bind:value={email}
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-700 bg-gray-800 text-white placeholder-gray-400 {isLogin ? 'rounded-t-md' : ''} focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Email address"
          />
        </div>
        <div>
          <label for="password" class="sr-only">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            required
            bind:value={password}
            class="appearance-none rounded-none relative block w-full px-3 py-2 border border-gray-700 bg-gray-800 text-white placeholder-gray-400 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 focus:z-10 sm:text-sm"
            placeholder="Password"
          />
        </div>
      </div>

      {#if form?.error}
        <div class="text-red-500 text-sm text-center">{form.error}</div>
      {/if}

      <div>
        <button
          type="submit"
          class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          {isLogin ? 'Sign in' : 'Create account'}
        </button>
      </div>

      <div class="text-center">
        <button
          type="button"
          class="text-indigo-400 hover:text-indigo-300"
          on:click={() => (isLogin = !isLogin)}
        >
          {isLogin ? 'Need an account? Sign up' : 'Already have an account? Sign in'}
        </button>
      </div>
    </form>
  </div>
</div> 