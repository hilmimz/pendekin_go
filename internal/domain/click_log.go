package domain

import "time"

type ClickLog struct {
	ID         int       `json:"id"`
	ClickedAt  time.Time `json:"clicked_at"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Referer    string    `json:"referrer"`
	ShortURLID int       `json:"short_url_id"`
}

type ClickLogRepository interface {
	Create(clickLog *ClickLog) error
}
