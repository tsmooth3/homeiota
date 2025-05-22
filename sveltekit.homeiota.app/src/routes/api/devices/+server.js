import { pool } from '$lib/db';

export async function GET() {
  try {
    // Fetch pump data
    const pumpResult = await pool.query(`
      SELECT 
        id,
        run_time,
        current,
        low_current,
        timestamp
      FROM pump_run_times
      ORDER BY timestamp DESC
      LIMIT 1
    `);
    // Fetch pump data
    const lowpumpResult = await pool.query(`
      WITH consecutive_readings AS (
        SELECT 
          t1.id,
          t1.timestamp as first_timestamp,
          t2.timestamp as second_timestamp,
          t1.current as first_current,
          t2.current as second_current
        FROM pump_run_times_critical t1
        JOIN pump_run_times_critical t2 
          ON t1.id <> t2.id 
          AND t2.timestamp > t1.timestamp 
          AND t2.timestamp <= t1.timestamp + interval '6 hours'
        WHERE t1.low_current = true AND t2.low_current = true
        ORDER BY t1.timestamp DESC
        LIMIT 1
      )
      SELECT 
        first_timestamp as timestamp,
        true as low_current
      FROM consecutive_readings
    `);

    // Fetch latest temperature for each location
    const tempResult = await pool.query(`
      WITH latest_temps AS (
        SELECT DISTINCT ON (location)
          id,
          value,
          location,
          timestamp
        FROM temperatures
        ORDER BY location, timestamp DESC
      )
      SELECT * FROM latest_temps
    `);

    // Fetch latest pump heartbeat
    const heartbeatResult = await pool.query(`
      SELECT 
        timestamp
      FROM device_heartbeats
      WHERE pump = true
      ORDER BY timestamp DESC
      LIMIT 1
    `);

    const now = new Date();
    const oneHourAgo = new Date(now - 60 * 60 * 1000);
    const tenMinutesAgo = new Date(now - 10 * 60 * 1000);

    const getStatus = (timestamp) => {
      const lastUpdate = new Date(timestamp);
      if (lastUpdate < oneHourAgo) return 'offline';
      if (lastUpdate < tenMinutesAgo) return 'warning';
      return 'online';
    };

    const getMinutesAgo = (timestamp) => {
      const lastUpdate = new Date(timestamp);
      return Math.floor((now - lastUpdate) / (60 * 1000));
    };
    const getHoursAgo = (timestamp) => {
      const lastUpdate = new Date(timestamp);
      return Math.floor((now - lastUpdate) / (60 * 60 * 1000));
    };

    // Transform data into device format
    const devices = [
      {
        id: 'pump',
        name: 'Well Pump',
        status: getStatus(heartbeatResult.rows[0]?.timestamp || pumpResult.rows[0]?.timestamp || new Date(0)),
        currentValue: pumpResult.rows[0]?.current || 0,
        lastHeartbeat: heartbeatResult.rows[0]?.timestamp || new Date(),
        lastUpdate: pumpResult.rows[0]?.timestamp || new Date(),
        lastLowUpdate: lowpumpResult.rows[0]?.timestamp || 0,
        minutesAgo: getMinutesAgo(pumpResult.rows[0]?.timestamp || new Date()),
        details: {
          runTime: pumpResult.rows[0]?.run_time || 0,
          lowCurrent: pumpResult.rows[0]?.low_current || false
        }
      },
      ...tempResult.rows.map(temp => ({
        id: `${temp.location}`,
        name: `Temperature - ${temp.location}`,
        status: getStatus(temp.timestamp),
        currentValue: temp.value,
        lastUpdate: temp.timestamp,
        minutesAgo: getMinutesAgo(temp.timestamp),
        details: {
          location: temp.location
        }
      }))
    ];

    return new Response(JSON.stringify({ devices }), {
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