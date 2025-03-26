package utils

const Symbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

const GetURLRegular = `SELECT original_url FROM urls WHERE short_url = $1`

const GetExistingURLRegular = `SELECT short_url FROM urls WHERE original_url = $1`

const GetURLsByUserID = "SELECT short_url, original_url FROM urls WHERE user_id = $1"

const SetURLRegular = `INSERT INTO urls (short_url, original_url, user_id) VALUES ($1, $2, $3) ON CONFLICT (original_url) DO NOTHING RETURNING short_url`
