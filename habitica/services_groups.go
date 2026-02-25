package habitica

import (
	"context"
	"fmt"
)

// GroupsService wraps group-related endpoints (party, guilds, tavern).
type GroupsService struct {
	client *Client
}

// GetGroup fetches a single group (GET /groups/:groupId).
func (s *GroupsService) GetGroup(ctx context.Context, id UUID) (*Group, error) {
	var g Group
	if err := s.client.doRequest(ctx, "GET", fmt.Sprintf("/groups/%s", id), nil, nil, &g); err != nil {
		return nil, err
	}
	return &g, nil
}

