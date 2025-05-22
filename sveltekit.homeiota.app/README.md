# Device Monitor

A SvelteKit application for monitoring device status and metrics in real-time.

## Features

- Real-time device status monitoring
- Interactive time series charts
- Responsive data tables
- Auto-updating data
- Mobile-friendly design

## Prerequisites

- Node.js 18 or later
- npm or yarn

## Setup

1. Install dependencies:
```bash
npm install
```

2. Start the development server:
```bash
npm run dev
```

3. Open [http://localhost:5173](http://localhost:5173) in your browser.

## Building for Production

```bash
npm run build
```

## Project Structure

- `src/routes/+page.svelte` - Main page component
- `src/routes/api/devices/+server.js` - API endpoint for device data
- `src/app.css` - Global styles
- `src/routes/+layout.svelte` - Main layout component

## Customization

To integrate with your own API:

1. Modify the `generateMockData` function in `src/routes/api/devices/+server.js`
2. Update the API endpoint URL in `src/routes/+page.svelte`
3. Adjust the data structure to match your API response

## Technologies Used

- SvelteKit
- TailwindCSS
- Chart.js
- Svelte Chart.js 