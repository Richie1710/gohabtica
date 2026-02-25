package habitica

// Webhook represents a registered Habitica webhook.
type Webhook struct {
	ID        UUID     `json:"id"`
	URL       string   `json:"url"`
	Enabled   bool     `json:"enabled"`
	Label     string   `json:"label"`
	Type      string   `json:"type"`
	Options   map[string]any `json:"options"`
	CreatedAt Timestamp `json:"createdAt"`
}

