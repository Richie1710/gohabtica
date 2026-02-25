package habitica

import (
	"context"
)

// TagsService wraps tag-related endpoints.
type TagsService struct {
	client *Client
}

// ListTags retrieves all tags of the user (GET /tags).
func (s *TagsService) ListTags(ctx context.Context) ([]*Tag, error) {
	var tags []*Tag
	if err := s.client.doRequest(ctx, "GET", "/tags", nil, nil, &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

