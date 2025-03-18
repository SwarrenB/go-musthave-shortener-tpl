package utils

const Symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const CreateTableQuery = `CREATE TABLE IF NOT EXISTS urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_url TEXT NOT NULL,
    original_url TEXT NOT NULL
)`

const GetURLRegular = `SELECT original_url FROM urls WHERE short_url = $1`

const SetURLRegular = `INSERT INTO urls (short_url, original_url) VALUES ($1, $2)`
