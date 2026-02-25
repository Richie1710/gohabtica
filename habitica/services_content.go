package habitica

import "context"

// ContentService wraps the static content endpoint.
type ContentService struct {
	client *Client
}

// GetContent retrieves the full static content payload (GET /content).
func (s *ContentService) GetContent(ctx context.Context) (*Content, error) {
	var c Content
	if err := s.client.doRequest(ctx, "GET", "/content", nil, nil, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

