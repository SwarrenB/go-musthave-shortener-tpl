package utils

const createTableQuery = `CREATE TABLE IF NOT EXISTS urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_url TEXT NOT NULL,
    original_url TEXT NOT NULL
)`

const getURLRegular = `SELECT original_url FROM urls WHERE short_url = $1`

const setURLRegular = `INSERT INTO urls (short_url, original_url) VALUES ($1, $2) ON CONFLICT (original_url) DO NOTHING`
