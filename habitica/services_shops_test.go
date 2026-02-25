package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShopsService_GetMarket(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/user/inventory/buy", r.URL.Path)

		resp := APIResponse[[]*ShopItem]{
			Success: true,
			Data: []*ShopItem{
				{Key: "potion", Text: "Health Potion"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	items, err := client.Shops.GetMarket(context.Background())
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, "potion", items[0].Key)
}

