package utils

const getURLRegular = `SELECT original_url FROM urls WHERE short_url = $1`

const setURLRegular = `INSERT INTO urls (short_url, original_url) VALUES ($1, $2)`
