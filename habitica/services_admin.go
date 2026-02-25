package habitica

import (
	"context"
	"fmt"
)

// AdminService wraps admin/moderator endpoints.
type AdminService struct {
	client *Client
}

// GetUserHistory retrieves the history of a user (GET /admin/user/:userId/history).
func (s *AdminService) GetUserHistory(ctx context.Context, userIdentifier string) ([]*UserHistoryEntry, error) {
	var hist []*UserHistoryEntry
	if err := s.client.doRequest(ctx, "GET", fmt.Sprintf("/admin/user/%s/history", userIdentifier), nil, nil, &hist); err != nil {
		return nil, err
	}
	return hist, nil
}

