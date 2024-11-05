package marshal

type URLRequest struct {
	OriginalURL string `json:"url"`
}

type URLResponse struct {
	ShortURL string `json:"result"`
}
