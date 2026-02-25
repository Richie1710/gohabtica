package habitica

import (
	"context"
	"net/url"
	"strconv"
)

// UserService groups endpoints around the currently authenticated user.
type UserService struct {
	client *Client
}

// GetCurrent fetches the currently authenticated user (GET /user).
func (s *UserService) GetCurrent(ctx context.Context) (*User, error) {
	var u User
	if err := s.client.doRequest(ctx, "GET", "/user", nil, nil, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetInbox fetches the user's inbox messages (GET /inbox/messages).
func (s *UserService) GetInbox(ctx context.Context, page int) (map[string]any, error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	var data map[string]any
	if err := s.client.doRequest(ctx, "GET", "/inbox/messages", q, nil, &data); err != nil {
		return nil, err
	}
	return data, nil
}

