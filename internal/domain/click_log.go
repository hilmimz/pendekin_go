package domain

type ClickLog struct {
	ID          int    `json:"id"`
	ClickedAt   string `json:"clicked_at"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
	Referrer    string `json:"referrer"`
	ShortLinkID int    `json:"short_link_id"`
}
