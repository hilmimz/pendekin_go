package domain

type ShortLink struct {
	ID          int    `json:"id"`
	OriginalURL string `json:"original_url"`
	ClickCount  int    `json:"click_count"`
	CreatedAt   string `json:"created_at"`
	UserID      int    `json:"user_id"`
	ExpiresIn   int    `json:"expires_in"`
	Alias       string `json:"alias"`
}
