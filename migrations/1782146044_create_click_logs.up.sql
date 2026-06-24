CREATE TABLE IF NOT EXISTS click_logs (
    id SERIAL PRIMARY KEY,
    clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT NOT NULL,
    referer TEXT,
    short_link_id INTEGER NOT NULL,
    FOREIGN KEY (short_link_id) REFERENCES short_links(id) ON DELETE CASCADE
);