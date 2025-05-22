<script>
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import DeviceDetailsTable from '$lib/components/DeviceDetailsTable.svelte';
  import TimeSeriesChart from '$lib/components/TimeSeriesChart.svelte';

  let device = null;
  let devicefiltered = null;
  let deviceCritical = null; // For pump run times critical
  let loading = true;
  let error = null;
  let isPump = false;
  let location = null;
  let currentPage = 1; // Track the current page for lazy loading
  const pageSize = 10; // Number of records to load per page

  onMount(async () => {
    await fetchDeviceData();
    if ($page.params.id === 'pump') {
      await fetchDeviceCriticalData();
    } else {
      loading = false;
    }
  });

  async function fetchDeviceData() {
    try {
      const response = await fetch(`/api/devices/${$page.params.id}?page=${currentPage}&pageSize=${pageSize}`);
      const data = await response.json();
      if (device === null) {
        device = data.device;
      } else {
        device = [...device, ...data.device]; // Append new data to existing data
      }
      isPump = $page.params.id === 'pump';
      // filter device for pump to show rows where current value < 1
      if (isPump) {
        devicefiltered = device.filter(d => d.current < 1);
        location = 'Well Pump Details'
      } else {
        location = device[0]?.location + ' Details';
      }
    } catch (e) {
      error = e.message;
    }
  }

  // Fetch all pump run times critical rows
  async function fetchDeviceCriticalData() {
    try {
      const response = await fetch('/api/devices/pump/critical');
      const data = await response.json();
      deviceCritical = data.deviceCritical;
      loading = false;
    } catch (e) {
      // Don't block page load on this error
      deviceCritical = null;
      loading = false;
    }
  }

  // Function to load more data when the user scrolls to the bottom
  function loadMore() {
    currentPage += 1; // Increment the page number
    fetchDeviceData(); // Fetch the next set of data
  }
</script>

<div class="space-y-8">
  <div class="flex items-center gap-4">
    <h1 class="text-2xl font-bold text-white">{location || 'Device Details'}</h1>
  </div>

  {#if loading}
    <div class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-12 w-12 border-4 border-blue-500 border-t-transparent"></div>
    </div>
  {:else if error}
    <div class="bg-red-900/50 border border-red-500 text-red-200 px-6 py-4 rounded-lg">
      <p>Error: {error}</p>
    </div>
  {:else if device}
    <div class="grid gap-8">
      <TimeSeriesChart data={device} {isPump}/>
      {#if isPump}
        {#if deviceCritical}
            <div>
              <h2 class="text-lg font-bold text-white mb-2">Pump Run Times Critical</h2>
              <DeviceDetailsTable device={deviceCritical} isPump={true} />
            </div>
          {/if}  
        {#if devicefiltered.length > 0}
        <div>
          <h2 class="text-lg font-bold text-white mb-2">Pump Run Times</h2>
          <DeviceDetailsTable device={devicefiltered} {isPump} />
          <h2 class="text-lg font-bold text-white mb-2">All Pump Readings</h2>
          <DeviceDetailsTable {device} {isPump} />
        </div>  
        {:else}
          <DeviceDetailsTable {device} {isPump} />
        {/if}
        
      {:else}
      <div>
        <h2 class="text-lg font-bold text-white mb-2">All {device[0]?.location} Readings</h2>
        <DeviceDetailsTable {device} {isPump} />
      </div>
      {/if}
      <button
        class="mt-4 inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        on:click={loadMore}
      >
        Load More
      </button>
    </div>
  {/if}
</div> 