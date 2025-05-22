<script>
  export let device;

  function getDaysAgo(timestamp) {
    const now = new Date();
    const lastUpdate = new Date(timestamp);
    const diffMs = now - lastUpdate;
    return Math.floor(diffMs / (1000 * 60 * 60 * 24));
  }
  function getHoursAgo(timestamp) {
    const now = new Date();
    const lastUpdate = new Date(timestamp);
    const diffMs = now - lastUpdate;
    return Math.floor(diffMs / (1000 * 60 * 60));
  }
  function getMinutesAgo(timestamp) {
    const now = new Date();
    const lastUpdate = new Date(timestamp);
    const diffMs = now - lastUpdate;
    return Math.floor(diffMs / (1000 * 60));
  }
  function getSecondsAgo(timestamp) {
    const now = new Date();
    const lastUpdate = new Date(timestamp);
    const diffMs = now - lastUpdate;
    return Math.floor(diffMs / (1000));
  }
</script>

<a href="/devices/{device.id}" class="flex-1 bg-gray-800 rounded-xl shadow-lg p-6 border border-gray-700 hover:border-gray-600 transition-colors cursor-pointer">
  <div class="flex items-center justify-between">
    <h3 class="text-lg font-semibold text-white">{device.name}</h3>
    <span class="px-3 py-1 rounded-full text-sm font-medium {device.status === 'online' ? 'bg-green-500/20 text-green-400' : 'bg-red-500/20 text-red-400'}">
      {device.status}
    </span>
  </div>
  <div class="mt-4 space-y-2">
    <div class="flex flex-col gap-1">
      <p class="text-gray-400 text-sm">Last Reading: {new Date(device.lastUpdate).toLocaleString()}</p>
      {#if getHoursAgo(device.lastUpdate) > 0}
        <p class="text-gray-400 text-sm">Updated: {getHoursAgo(device.lastUpdate)} hours ago</p>
      {:else if getMinutesAgo(device.lastUpdate) > 0}
        <p class="text-gray-400 text-sm">Updated: {getMinutesAgo(device.lastUpdate)} minutes ago</p>
      {:else}
        <p class="text-gray-400 text-sm">Updated: {getSecondsAgo(device.lastUpdate)} seconds ago</p>
      {/if}
      {#if device.name.toLowerCase().includes('well pump')}
        {#if getHoursAgo(device.lastHeartbeat) > 0}
          <p class="text-gray-400 text-sm">Heartbeat: {getHoursAgo(device.lastHeartbeat)} hours ago</p>
        {:else if  getMinutesAgo(device.lastHeartbeat) > 0}
          <p class="text-gray-400 text-sm">Heartbeat: {getMinutesAgo(device.lastHeartbeat)} minutes ago</p>
        {:else}
          <p class="text-gray-400 text-sm">Heartbeat: {getSecondsAgo(device.lastHeartbeat)} seconds ago</p>
        {/if}
        {#if device.lastLowUpdate != 0}
          <p class="text-gray-400 text-sm">Last Dry: {new Date(device.lastLowUpdate).toLocaleString()}</p>
          {#if getDaysAgo(device.lastLowUpdate) > 0}
            <p class="text-gray-400 text-sm">Last Dry: {getDaysAgo(device.lastLowUpdate)} days ago</p>
          {:else if getHoursAgo(device.lastLowUpdate) > 0}
            <p class="text-gray-400 text-sm">Last Dry: {getHoursAgo(device.lastLowUpdate)} hours ago</p>
          {:else if getMinutesAgo(device.lastLowUpdate) > 0}
            <p class="text-gray-400 text-sm">Last Dry: {getMinutesAgo(device.lastLowUpdate)} minutes ago</p>
          {:else}
            <p class="text-gray-400 text-sm">Last Dry: {getSecondsAgo(device.lastLowUpdate)} seconds ago</p>
          {/if}
        {/if}
      {/if}

    </div>
    <div class="flex items-center gap-2">
      <div class="w-2 h-2 rounded-full {device.status === 'online' ? 'bg-green-500' : 'bg-red-500'}"></div>
      <p class="text-gray-300">
        {#if device.name.toLowerCase().includes('well pump')}
          Run Time: {device.details.runTime} seconds
        {:else}
          Current Value: {device.currentValue}
        {/if}
      </p>
    </div>
  </div>
</a> 