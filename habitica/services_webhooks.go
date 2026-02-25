package habitica

import (
	"context"
)

// WebhooksService wraps webhook-related endpoints.
type WebhooksService struct {
	client *Client
}

// ListWebhooks retrieves all registered webhooks of the user (GET /user/webhook).
func (s *WebhooksService) ListWebhooks(ctx context.Context) ([]*Webhook, error) {
	var hooks []*Webhook
	if err := s.client.doRequest(ctx, "GET", "/user/webhook", nil, nil, &hooks); err != nil {
		return nil, err
	}
	return hooks, nil
}

