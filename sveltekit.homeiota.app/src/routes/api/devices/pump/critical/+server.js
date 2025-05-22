import { pool } from '$lib/db';

export async function GET() {
  try {
    const result = await pool.query(`
      SELECT 
        id,
        run_time,
        current,
        low_current,
        timestamp
      FROM pump_run_times_critical
      ORDER BY timestamp DESC
    `);
    return new Response(
      JSON.stringify({ deviceCritical: result.rows }),
      { headers: { 'Content-Type': 'application/json' } }
    );
  } catch (error) {
    console.error('Database error:', error);
    return new Response(
      JSON.stringify({ error: 'Failed to fetch pump run times critical' }),
      { status: 500, headers: { 'Content-Type': 'application/json' } }
    );
  }
} 