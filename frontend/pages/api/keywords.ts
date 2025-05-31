import type { NextApiRequest, NextApiResponse } from 'next';
import { Pool } from 'pg';

const pool = new Pool({
  connectionString: process.env.DATABASE_URL,
});

// Ensure table exists
async function ensureTable() {
  await pool.query(`
    CREATE TABLE IF NOT EXISTS keywords (
      id SERIAL PRIMARY KEY,
      keyword TEXT NOT NULL,
      response TEXT NOT NULL
    )
  `);
}

export default async function handler(req: NextApiRequest, res: NextApiResponse) {
  await ensureTable();

  if (req.method === 'GET') {
    // Fetch all keyword-response pairs
    const result = await pool.query('SELECT keyword, response FROM keywords ORDER BY id ASC');
    return res.status(200).json(result.rows);
  }

  if (req.method === 'POST') {
    const pairs = Array.isArray(req.body) ? req.body : [];
    if (!pairs.every(p => typeof p.keyword === 'string' && typeof p.response === 'string')) {
      return res.status(400).json({ error: 'Invalid data format' });
    }
    // Clear and insert new pairs
    await pool.query('TRUNCATE TABLE keywords');
    for (const { keyword, response } of pairs) {
      await pool.query('INSERT INTO keywords (keyword, response) VALUES ($1, $2)', [keyword, response]);
    }
    return res.status(200).json({ ok: true });
  }

  res.setHeader('Allow', ['GET', 'POST']);
  res.status(405).end(`Method ${req.method} Not Allowed`);
}
