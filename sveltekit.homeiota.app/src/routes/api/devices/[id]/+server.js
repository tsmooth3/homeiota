import { json } from '@sveltejs/kit';
import { pool } from '$lib/db';

export async function GET({ params }) {
  try {
    const { id } = params;
    let result;
    if (id === 'pump') {
      result = await pool.query(
        'SELECT * FROM pump_run_times ORDER BY timestamp DESC'
      );
    } else {
      result = await pool.query(
        'SELECT * FROM temperatures WHERE location = $1 ORDER BY timestamp DESC LIMIT 10000', 
        [id]
      );
    }

    if (result.rows.length === 0) {
      return new Response('Device not found', { status: 404 });
    }

    const device = result.rows;
    return json({ device });
  } catch (error) {
    console.error('Error fetching device:', error);
    return new Response('Internal Server Error', { status: 500 });
  }
} 