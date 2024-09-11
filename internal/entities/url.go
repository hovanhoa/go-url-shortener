package entities

type URL struct {
	ID        int64  `json:"id"`
	SortURL   string `json:"sort_url"`
	LongURL   string `json:"long_url"`
	CreatedAt int64  `json:"created_at"`
}
