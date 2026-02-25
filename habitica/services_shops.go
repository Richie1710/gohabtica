package habitica

import "context"

// ShopsService wraps shop/market-related endpoints.
type ShopsService struct {
	client *Client
}

// GetMarket retrieves market items (GET /user/inventory/buy).
func (s *ShopsService) GetMarket(ctx context.Context) ([]*ShopItem, error) {
	var items []*ShopItem
	if err := s.client.doRequest(ctx, "GET", "/user/inventory/buy", nil, nil, &items); err != nil {
		return nil, err
	}
	return items, nil
}

