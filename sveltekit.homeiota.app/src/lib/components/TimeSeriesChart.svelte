<script>
  import { onMount, onDestroy } from 'svelte';
  import * as d3 from 'd3';

  export let data;
  export let isPump = false;
  let TEMP_THRESHOLD = 10.0; // Default value for temperature warning threshold
  if (data[0]['location'] === 'freezer') {
    TEMP_THRESHOLD = 10.0; // Default value for temperature warning threshold;
  } else {
    TEMP_THRESHOLD = 75.0; // Default value for temperature warning threshold;
  }
  const DAYS_TO_SHOW = 3;

  let chartContainer;
  let resizeObserver;

  function drawChart() {
    // Clear previous chart
    d3.select(chartContainer).selectAll('*').remove();

    // Sort data by timestamp in ascending order
    const sortedData = [...data].sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

    // Filter data to last 3 days
    const threeDaysAgo = new Date();
    threeDaysAgo.setDate(threeDaysAgo.getDate() - DAYS_TO_SHOW);
    const filteredData = sortedData.filter(d => new Date(d.timestamp) >= threeDaysAgo);

    let chartData;
    if (isPump) {
      // For pumps, use raw data but only show where current value = 0
      chartData = filteredData.filter(d => d.current < 1).map(d => ({ 
        x: new Date(d.timestamp),
        y: d.run_time
      }));
    } else {
      // For temperature sensors, group by 5-minute intervals and calculate averages
      const groupedData = new Map();
      filteredData.forEach(d => {
        const date = new Date(d.timestamp);
        // Round down to nearest 5 minutes
        const roundedDate = new Date(Math.floor(date.getTime() / (5 * 60 * 1000)) * (5 * 60 * 1000));
        const key = roundedDate.getTime();
        
        if (!groupedData.has(key)) {
          groupedData.set(key, {
            timestamp: roundedDate,
            values: [],
            count: 0
          });
        }
        
        const group = groupedData.get(key);
        group.values.push(d.value);
        group.count++;
      });

      // Calculate averages for each group
      chartData = Array.from(groupedData.values()).map(group => ({
        x: group.timestamp,
        y: group.values.reduce((sum, val) => sum + val, 0) / group.count
      }));
    }

    // Get container dimensions
    const width = chartContainer.clientWidth;
    const height = chartContainer.clientHeight;
    const margin = { top: 20, right: 20, bottom: 50, left: 40 };
    const innerWidth = width - margin.left - margin.right;
    const innerHeight = height - margin.top - margin.bottom;

    // Create SVG
    const svg = d3.select(chartContainer)
      .append('svg')
      .attr('width', width)
      .attr('height', height);

    const g = svg.append('g')
      .attr('transform', `translate(${margin.left},${margin.top})`);

    // Create scales
    const xScale = d3.scaleTime()
      .domain(d3.extent(chartData, d => d.x))
      .range([0, innerWidth]);

    const yScale = d3.scaleLinear()
      .domain(d3.extent(chartData, d => d.y))
      .range([innerHeight, 0]);

    // Create line generator with linear curve
    const line = d3.line()
      .x(d => xScale(d.x))
      .y(d => yScale(d.y))
      .curve(d3.curveLinear);

    // Add grid lines
    g.append('g')
      .attr('class', 'grid')
      .attr('transform', `translate(0,${innerHeight})`)
      .call(d3.axisBottom(xScale)
        .tickSize(-innerHeight)
        .tickFormat(''))
      .attr('stroke', '#374151')
      .attr('stroke-opacity', 0.5);

    // Add x-axis
    g.append('g')
      .attr('transform', `translate(0,${innerHeight})`)
      .call(d3.axisBottom(xScale)
        .tickFormat(d3.timeFormat('%m/%d %H:%M')))
      .attr('color', '#9CA3AF')
      .attr('font-size', '10px')
      .selectAll('text')
      .attr('transform', 'rotate(-45)')
      .attr('text-anchor', 'end')
      .attr('dx', '-0.5em')
      .attr('dy', '0.5em');

    // Add y-axis
    g.append('g')
      .call(d3.axisLeft(yScale)
        .tickFormat(d => isPump ? `${d} sec` : `${d}°F`))
      .attr('color', '#9CA3AF')
      .attr('font-size', '10px');

    if (!isPump) {
      // Add threshold line for temperature
      g.append('line')
        .attr('x1', 0)
        .attr('x2', innerWidth)
        .attr('y1', yScale(TEMP_THRESHOLD))
        .attr('y2', yScale(TEMP_THRESHOLD))
        .attr('stroke', '#EF4444')
        .attr('stroke-width', 1)
        .attr('stroke-dasharray', '4,4');

      // Add threshold label
      g.append('text')
        .attr('x', innerWidth - 5)
        .attr('y', yScale(TEMP_THRESHOLD) - 5)
        .attr('text-anchor', 'end')
        .attr('fill', '#EF4444')
        .attr('font-size', '10px')
        .text(`${TEMP_THRESHOLD}°F Threshold`);

      // Add area above threshold with linear curve
      const area = d3.area()
        .x(d => xScale(d.x))
        .y0(yScale(TEMP_THRESHOLD))
        .y1(d => yScale(d.y))
        .curve(d3.curveLinear);

      g.append('path')
        .datum(chartData.filter(d => d.y > TEMP_THRESHOLD))
        .attr('fill', '#EF4444')
        .attr('fill-opacity', 0.1)
        .attr('d', area);
    } else {
      // add threshold line for pump at 60 seconds
      g.append('line')
        .attr('x1', 0)
        .attr('x2', innerWidth)
        .attr('y1', yScale(60))
        .attr('y2', yScale(60))
        .attr('stroke', '#60A5FA')
        .attr('stroke-width', 1)
        .attr('stroke-dasharray', '4,4');

      // add threshold label
      g.append('text')
        .attr('x', innerWidth - 5)
        .attr('y', yScale(60) - 5)
        .attr('text-anchor', 'end')
        .attr('fill', '#60A5FA')
        .attr('font-size', '10px')
        .text('60 sec Threshold');
        
    }

    // Add line
    g.append('path')
      .datum(chartData)
      .attr('fill', 'none')
      .attr('stroke', '#60A5FA')
      .attr('stroke-width', 2)
      .attr('d', line);
  }

  onMount(() => {
    // Initial draw
    drawChart();

    // Set up resize observer
    resizeObserver = new ResizeObserver(() => {
      drawChart();
    });

    resizeObserver.observe(chartContainer);
  });

  onDestroy(() => {
    if (resizeObserver) {
      resizeObserver.disconnect();
    }
  });
</script>

<div class="bg-gray-800 rounded-xl shadow-lg p-4 sm:p-6 border border-gray-700 w-full max-w-full min-w-0 overflow-x-auto">
  <h3 class="text-lg font-semibold text-white mb-4">
    Historical Data {#if !isPump}(5-minute averages){/if}
  </h3>
  <div class="h-80 w-full min-w-0" bind:this={chartContainer}></div>
  <!-- Optionally, for aspect ratio-based scaling, replace the above div with:
  <div class='relative w-full aspect-[2/1] min-w-0' bind:this={chartContainer}></div>
  -->
</div> 