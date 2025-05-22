<script>
  export let device;
  export let isPump;
  // Pagination state
  let currentPage = 1;
  const pageSize = 10;

  function prevPage() {
    if (currentPage > 1) currentPage--;
  }
  function nextPage() {
    if (currentPage < totalPages) currentPage++;
  }

  // Collapsible state
  let collapsed = false;

  // Filtering state
  let filterValue = '';
  let filterType = 'under'; // for temperature: 'over' or 'under'
  let filterMin = '';
  let filterMax = '';

  // Filtered data
  $: filteredDevice = device ? device.filter(detail => {
    if (isPump) {
      const min = filterMin === '' ? -Infinity : Number(filterMin);
      const max = filterMax === '' ? Infinity : Number(filterMax);
      return detail.current >= min && detail.current <= max;
    } else {
      if (filterValue === '' ||filterValue === null || filterValue === undefined || filterValue === 0 || isNaN(Number(filterValue))) return true;
      const num = Number(filterValue);
      return filterType === 'under' ? detail.value < num : detail.value > num;
    }
  }) : [];
  $: totalPages = filteredDevice ? Math.ceil(filteredDevice.length / pageSize) : 1;
  $: pagedDevice = filteredDevice ? filteredDevice.slice((currentPage - 1) * pageSize, currentPage * pageSize) : [];
</script>

<div class="bg-gray-800 rounded-xl shadow-lg overflow-hidden border border-gray-700 max-w-full sm:max-w-screen-sm mx-auto mb-4">
  <div class="flex items-center justify-between p-4 bg-gray-900 border-b border-gray-700">
    <button on:click={() => collapsed = !collapsed} class="px-3 py-1 rounded bg-gray-700 text-gray-200 flex items-center justify-center" aria-label={collapsed ? 'Show Table' : 'Hide Table'}>
      {#if collapsed}
        <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
        </svg>
      {:else}
        <svg class="w-5 h-5" fill="none" stroke="currentColor" stroke-width="2" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" d="M5 9l7 7 7-7" />
        </svg>
      {/if}
    </button>
    {#if isPump}
      <div class="flex items-center gap-2">
        <label for="current-min" class="text-gray-400 text-sm">Current</label>
        <input id="current-min" type="number" bind:value={filterMin} class="w-20 px-2 py-1 rounded bg-gray-700 text-gray-200 border border-gray-600" placeholder="Min" min="0" step="any" />
        <span class="text-gray-400 text-sm">to</span>
        <input id="current-max" type="number" bind:value={filterMax} class="w-20 px-2 py-1 rounded bg-gray-700 text-gray-200 border border-gray-600" placeholder="Max" min="0" step="any" />
      </div>
    {:else}
      <div class="flex items-center gap-2">
        <label for="value-filter" class="text-gray-400 text-sm">Value</label>
        <select bind:value={filterType} class="px-1 py-1 rounded bg-gray-700 text-gray-200 border border-gray-600">
          <option value="under">&lt;</option>
          <option value="over">&gt;</option>
        </select>
        <input id="value-filter" type="number" bind:value={filterValue} class="w-20 px-2 py-1 rounded bg-gray-700 text-gray-200 border border-gray-600" placeholder="Value" step="any" />
      </div>
    {/if}
  </div>
  {#if !collapsed}
    <div class="overflow-x-auto">
      <table class="min-w-full w-full divide-y divide-gray-700">
        <thead class="bg-gray-800">
          <tr>
            <th class="px-2 sm:px-6 py-4 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Timestamp</th>
            {#if isPump}
              <th class="px-2 sm:px-6 py-4 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Run Time</th>
              <th class="px-2 sm:px-6 py-4 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Current</th>
            {:else}
              <th class="px-2 sm:px-6 py-4 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">Value</th>
            {/if}
          </tr>
        </thead>
        <tbody class="bg-gray-800 divide-y divide-gray-700">
          {#if isPump}
            {#each pagedDevice as detail}
            <tr class="hover:bg-gray-700/50 transition-colors">
              <td class="px-2 sm:px-6 py-4 text-sm text-gray-300">
                {new Date(detail.timestamp).toLocaleString()}
              </td>
              <td class="px-2 sm:px-6 py-4 text-sm text-gray-300">
                {detail.run_time} seconds
              </td>
              <td class="px-2 sm:px-6 py-4 text-sm text-gray-300">
                {detail.current}
              </td>
            </tr>
            {/each}
          {:else}
            {#each pagedDevice as detail}
            <tr class="hover:bg-gray-700/50 transition-colors">
              <td class="px-2 sm:px-6 py-4 text-sm text-gray-300">
                {new Date(detail.timestamp).toLocaleString()}
              </td>
              <td class="px-2 sm:px-6 py-4 text-sm text-gray-300">
                {detail.value}
              </td>
            </tr>
            {/each}
          {/if}
        </tbody>
      </table>
    </div>
    {#if filteredDevice && filteredDevice.length > pageSize}
      <div class="flex justify-between items-center p-4 bg-gray-900 border-t border-gray-700">
        <button on:click={prevPage} class="px-3 py-1 rounded bg-gray-700 text-gray-200 disabled:opacity-50" disabled={currentPage === 1}>&larr; Prev</button>
        <span class="text-gray-400">Page {currentPage} of {totalPages}</span>
        <button on:click={nextPage} class="px-3 py-1 rounded bg-gray-700 text-gray-200 disabled:opacity-50" disabled={currentPage === totalPages}>Next &rarr;</button>
      </div>
    {/if}
  {/if}
</div> 