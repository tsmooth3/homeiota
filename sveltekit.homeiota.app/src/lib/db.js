import pg from 'pg';
const { Pool } = pg;

// Create a new pool using environment variables
export const pool = new Pool({
  user: process.env.POSTGRES_USER,
  host: process.env.POSTGRES_HOST,
  database: process.env.POSTGRES_DB,
  password: process.env.POSTGRES_PASSWORD,
  port: process.env.POSTGRES_PORT || 5432,
});


// Test the connection
// pool.on('connect', () => {
//   console.log('Connected to PostgreSQL database');
// });

pool.on('error', (err) => {
  console.error('Unexpected error on idle client', err);
  process.exit(-1);
});