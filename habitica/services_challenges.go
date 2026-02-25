package habitica

import (
	"context"
	"fmt"
)

// ChallengesService wraps challenge-related endpoints.
type ChallengesService struct {
	client *Client
}

// GetChallenge fetches a single challenge (GET /challenges/:challengeId).
func (s *ChallengesService) GetChallenge(ctx context.Context, id UUID) (*Challenge, error) {
	var c Challenge
	if err := s.client.doRequest(ctx, "GET", fmt.Sprintf("/challenges/%s", id), nil, nil, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

