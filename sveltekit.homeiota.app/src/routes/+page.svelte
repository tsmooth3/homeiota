<script>
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { browser } from '$app/environment';
  import StatusCard from '$lib/components/StatusCard.svelte';

  let devices = [];
  let loading = true;
  let error = null;

  const fetchData = async () => {
    try {
      const response = await fetch('/api/devices', { cache: 'no-store' });
      if (!response.ok) throw new Error('Server error');
      const data = await response.json();
      devices = data.devices;
      error = null; // Clear error on success
      loading = false;
    } catch (e) {
      error = e.message || 'Failed to fetch';
      // Don't set loading to false here, so spinner can keep showing if you want
      // Optionally, you can set loading = false if you want to show the error box
    }
  };

  // Initial data fetch
  onMount(fetchData);

  // Refresh data when the page store changes (navigation)
  $: if ($page && browser) {
    fetchData();
  }

  // Poll for updates every 10 seconds
  let pollInterval;
  onMount(() => {
    pollInterval = setInterval(fetchData, 10000);
    return () => {
      if (pollInterval) clearInterval(pollInterval);
    };
  });

  
</script>

<div class="space-y-8">
  {#if loading}
    <div class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-blue-500 border-t-transparent"></div>
    </div>
  {:else if error}
    <div class="bg-red-900/50 border border-red-500 text-red-200 px-6 py-4 rounded-lg">
      <p>Error: {error}</p>
    </div>
  {:else}
    <!-- Status Cards -->
    <div class="flex flex-col md:flex-row gap-6">
      {#each devices as device}
        <StatusCard {device} />
      {/each}
    </div>
  {/if}
</div> 