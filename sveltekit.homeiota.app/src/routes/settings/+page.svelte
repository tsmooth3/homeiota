<script lang="ts">
  import { enhance } from '$app/forms';
  export let data;
  let { user, alertPreferences } = data;

  let name = user?.name ?? '';
  let email = user?.email ?? '';
  let gotifyToken = user?.gotifyToken || '';
  let showToken = false;
  let testStatus = '';
  let showAddAlertModal = false;
  let availableLocations = [];
  let newAlert = { location: '', threshold: 0, offlineThreshold: 0 };
  let saveStatus: 'idle' | 'saving' | 'success' | 'error' = 'idle';
  let saveMessage = '';

  // Only show devices with alert preferences for the signed-in user
  type uiAlertPreference = {
    name: string;
    threshold: number;
    enabled: boolean;
    offlineThreshold?: number;
  };
  $: uiAlertPreferences = alertPreferences
    ? alertPreferences.map(pref => ({
        name: pref.location,
        threshold: pref.threshold,
        enabled: pref.enabled,
        offlineThreshold: pref.offlineThreshold ?? 0
      }))
    : [];


  // Function to add a new alert
  async function addAlert() {
    try {
      const response = await fetch('/api/temperature-devices');
      const data = await response.json();
      if (response.ok) {
        // Assuming the response includes a list of devices with their locations and suggested thresholds
        const devices = data.devices;
        // You can implement a modal or dropdown here to select a device
        // For simplicity, we'll just add the first device in the list
        const newDevice = devices[0];
        uiAlertPreferences.push({
          name: newDevice.location,
          threshold: newDevice.suggestedThreshold, // Use the suggested threshold
          enabled: true, // Enable the alert by default
          offlineThreshold: newDevice.suggestedOfflineThreshold ?? 0
        });
      } else {
        console.error('Failed to fetch temperature devices:', data.error);
      }
    } catch (error) {
      console.error('Error adding alert:', error);
    }
  }

  // Accepts either a uiAlertPreference or 'offline' string
  async function testNotification(sensorOrType: uiAlertPreference | 'offline') {
    if (sensorOrType === 'offline') {
      try {
        testStatus = 'Sending test offline notification...';
        const response = await fetch('/api/test-notifications', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            type: 'offline',
            deviceId: 'test-device' // Replace with real device ID if needed
          })
        });
        const data = await response.json();
        if (response.ok) {
          testStatus = 'Test offline notification sent successfully!';
        } else {
          testStatus = `Error: ${data.error}`;
        }
      } catch (error) {
        testStatus = 'Failed to send test offline notification';
        console.error('Test offline notification error:', error);
      }
    } else {
      try {
        testStatus = `Sending test notification for ${sensorOrType.name}...`;
        const response = await fetch('/api/test-notifications', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            type: 'temperature',
            deviceId: sensorOrType.name,
            threshold: sensorOrType.threshold
          })
        });
        const data = await response.json();
        if (response.ok) {
          testStatus = `Test notification for ${sensorOrType.name} sent successfully!`;
        } else {
          testStatus = `Error: ${data.error}`;
        }
      } catch (error) {
        testStatus = 'Failed to send test notification';
        console.error('Test notification error:', error);
      }
    }
  }

  async function testOfflineThreshold(sensor: uiAlertPreference) {
    try {
      testStatus = `Sending offline test for ${sensor.name}...`;
      const response = await fetch('/api/test-notifications', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          type: 'offline',
          deviceId: sensor.name,
          threshold: sensor.offlineThreshold
        })
      });
      const data = await response.json();
      if (response.ok) {
        testStatus = `Offline test for ${sensor.name} sent successfully!`;
      } else {
        testStatus = `Error: ${data.error}`;
      }
    } catch (error) {
      testStatus = 'Failed to send offline test notification';
      console.error('Test offline notification error:', error);
    }
  }

  function handleSubmit() {
    return async ({ result }) => {
      if (result.type === 'success') {
        saveStatus = 'success';
        saveMessage = 'Settings saved successfully!';
        setTimeout(() => { saveStatus = 'idle'; saveMessage = ''; }, 3000);
        // Reload the page to get updated alertPreferences from the server
        window.location.reload();
        // Or, if you want to use SvelteKit navigation:
        // await goto(window.location.pathname, { replaceState: true });
      } else if (result.type === 'failure') {
        saveStatus = 'error';
        saveMessage = result.data?.error || 'Failed to save settings.';
        setTimeout(() => { saveStatus = 'idle'; saveMessage = ''; }, 5000);
      }
    };
  }

  async function openAddAlertModal() {
    showAddAlertModal = true;
    // Fetch available locations, excluding those already in uiAlertPreferences
    const response = await fetch('/api/temperature-devices');
    const data = await response.json();
    data.devices.push({location: "wellpump", suggestedThreshold: 0, suggestedOfflineThreshold: 0});
    console.log(data);
    const usedLocations = new Set(uiAlertPreferences.map(s => s.name));
    availableLocations = data.devices
      .map(d => d.location)
      .filter(loc => !usedLocations.has(loc));
    newAlert = {
      location: availableLocations[0] || '',
      threshold: 0,
      offlineThreshold: 0
    };
  }

  function closeAddAlertModal() {
    showAddAlertModal = false;
  }

  function addNewAlert() {
    if (!newAlert.location) return;
    uiAlertPreferences = [
      ...uiAlertPreferences,
      {
        name: newAlert.location,
        threshold: newAlert.threshold,
        enabled: true,
        offlineThreshold: newAlert.offlineThreshold
      }
    ];
    showAddAlertModal = false;
  }
</script>

<div class="max-w-2xl mx-auto">
  <h1 class="text-2xl font-bold mb-6">Settings</h1>

  {#if saveStatus === 'success'}
    <div class="mb-4 p-2 rounded bg-green-900/50 text-green-200 text-center">{saveMessage}</div>
  {:else if saveStatus === 'error'}
    <div class="mb-4 p-2 rounded bg-red-900/50 text-red-200 text-center">{saveMessage}</div>
  {/if}

  <form method="POST" use:enhance={handleSubmit} class="space-y-6">
    <input type="hidden" name="uiAlertPreferences" value={JSON.stringify(uiAlertPreferences)} />
    <div class="bg-gray-800 shadow rounded-lg p-6">
      <h2 class="text-lg font-medium mb-4">Profile</h2>
      <div class="space-y-4">
        <div>
          <label for="name" class="block text-sm font-medium text-gray-300">Name</label>
          <input
            type="text"
            id="name"
            name="name"
            bind:value={name}
            class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
        </div>
        <div>
          <label for="email" class="block text-sm font-medium text-gray-300">Email</label>
          <input
            type="email"
            id="email"
            name="email"
            bind:value={email}
            class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
          />
        </div>
        <div>
          <label for="gotifyToken" class="block text-sm font-medium text-gray-300">Gotify Client Token</label>
          <div class="relative">
            <input
              type="password"
              id="gotifyToken"
              name="gotifyToken"
              bind:value={gotifyToken}
              placeholder="Enter your Gotify client token"
              class="mt-1 block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm pr-10"
            />
            <button
              type="button"
              class="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-400 hover:text-gray-300"
              on:click={() => showToken = !showToken}
            >
              {#if showToken}
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
                  <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
                </svg>
              {:else}
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M3.707 2.293a1 1 0 00-1.414 1.414l14 14a1 1 0 001.414-1.414l-1.473-1.473A10.014 10.014 0 0019.542 10C18.268 5.943 14.478 3 10 3a9.958 9.958 0 00-4.512 1.074l-1.78-1.781zm4.261 4.26l1.514 1.515a2.003 2.003 0 012.45 2.45l1.514 1.514a4 4 0 00-5.478-5.478z" clip-rule="evenodd" />
                  <path d="M12.454 16.697L9.75 13.992a4 4 0 01-3.742-3.741L2.335 6.578A9.98 9.98 0 00.458 10c1.274 4.057 5.065 7 9.542 7 .847 0 1.669-.105 2.454-.303z" />
                </svg>
              {/if}
            </button>
          </div>
          {#if showToken}
            <div class="mt-2 p-2 bg-gray-700 rounded text-sm text-gray-300 break-all">
              {gotifyToken || 'No token set'}
            </div>
          {/if}
          <p class="mt-1 text-sm text-gray-400">
            This token will be used to send notifications to your Gotify client. You can find your client token in your Gotify dashboard.
          </p>
        </div>
      </div>
    </div>

    <div class="bg-gray-800 shadow rounded-lg p-6">
      <h2 class="text-lg font-medium mb-4">Alert Preferences</h2>
      <button
        type="button"
        class="mb-4 inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        on:click={openAddAlertModal}
      >
        Add Alert
      </button>
      {#if showAddAlertModal}
        <div class="fixed inset-0 z-50 bg-black/60 backdrop-blur-sm flex items-center justify-center">
          <div class="absolute bg-gray-900 rounded-lg shadow-lg w-full max-w-xs sm:max-w-sm  flex flex-col">
            <h3 class="text-base font-bold mb-2 text-white break-words">Add Alert Preference</h3>
            <div class="mb-3 w-full">
              <label for="add-alert-location" class="block text-xs font-medium text-gray-300 mb-1">Location</label>
              <select id="add-alert-location" bind:value={newAlert.location} class="block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 text-sm py-2 px-2 break-words">
                {#each availableLocations as loc}
                  <option value={loc}>{loc}</option>
                {/each}
              </select>
            </div>
            <div class="mb-3 w-full">
              <label for="add-alert-threshold" class="block text-xs font-medium text-gray-300 mb-1">Threshold (°F/Amps)</label>
              <input id="add-alert-threshold" type="number" bind:value={newAlert.threshold} class="block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 text-sm py-2 px-2 break-words" />
            </div>
            <div class="mb-3 w-full">
              <label for="add-alert-offline-threshold" class="block text-xs font-medium text-gray-300 mb-1">Offline Threshold (min)</label>
              <input id="add-alert-offline-threshold" type="number" min="0" bind:value={newAlert.offlineThreshold} class="block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 text-sm py-2 px-2 break-words" />
            </div>
            <div class="flex flex-col sm:flex-row justify-end gap-2 w-full mt-2">
              <button type="button" class="w-full sm:w-auto px-3 py-2 rounded bg-gray-700 text-gray-200 text-sm font-medium" on:click={closeAddAlertModal}>Cancel</button>
              <button type="button" class="w-full sm:w-auto px-3 py-2 rounded bg-indigo-600 text-white text-sm font-medium" on:click={addNewAlert}>Add</button>
            </div>
          </div>
        </div>
      {/if}
      <div class="max-w-[600px] w-full mx-auto overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-700">
          <thead>
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Device</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Threshold (°F/Amps)</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Enable Alerts</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-300 uppercase tracking-wider">Offline Threshold (min)</th>
            </tr>
          </thead>
          <tbody class="bg-gray-800 divide-y divide-gray-700">
            {#each uiAlertPreferences as sensor}
              <tr>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-300">{sensor.name}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    type="number"
                    bind:value={sensor.threshold}
                    class="block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                  />
                  <button
                    type="button"
                    class="w-full mt-2 inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                    on:click={() => testNotification(sensor)}
                  >
                    Test Alert
                  </button>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    type="checkbox"
                    bind:checked={sensor.enabled}
                    class="rounded bg-gray-700 border-gray-600 text-indigo-600 focus:ring-indigo-500"
                  />
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <input
                    type="number"
                    min="0"
                    bind:value={sensor.offlineThreshold}
                    class="block w-full rounded-md bg-gray-700 border-gray-600 text-white shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
                  />
                  <button
                    type="button"
                    class="w-full mt-2 inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md text-yellow-700 bg-yellow-100 hover:bg-yellow-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500"
                    on:click={() => testOfflineThreshold(sensor)}
                  >
                    Test Offline Threshold
                  </button>
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>

    
    <div class="flex justify-end">
      <button
        type="submit"
        class="inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        disabled={saveStatus === 'saving'}
      >
        {saveStatus === 'saving' ? 'Saving...' : 'Save Changes'}
      </button>
    </div>
  </form>
</div> 