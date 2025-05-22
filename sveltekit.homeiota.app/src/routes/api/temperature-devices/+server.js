import { pool } from '$lib/db';

export async function GET() {
  try {
    // Fetch the latest temperature for each location and calculate the suggested threshold
    const result = await pool.query(`
      WITH latest_temps AS (
        SELECT DISTINCT ON (location)
          location,
          value,
          timestamp
        FROM temperatures
        ORDER BY location, timestamp DESC
      )
      SELECT 
        location,
        MAX(value) as suggestedThreshold
      FROM latest_temps
      GROUP BY location
    `);

    return new Response(JSON.stringify({ devices: result.rows }), {
      headers: {
        'Content-Type': 'application/json'
      }
    });
  } catch (error) {
    console.error('Database error:', error);
    return new Response(JSON.stringify({ error: 'Failed to fetch data' }), {
      status: 500,
      headers: {
        'Content-Type': 'application/json'
      }
    });
  }
} 