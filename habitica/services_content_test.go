package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContentService_GetContent(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/content", r.URL.Path)

		resp := APIResponse[Content]{
			Success: true,
			Data: Content{
				Items: map[string]any{
					"weapon_warrior_0": map[string]any{"text": "Training Sword"},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	c, err := client.Content.GetContent(context.Background())
	require.NoError(t, err)
	require.NotNil(t, c)
	require.Contains(t, c.Items, "weapon_warrior_0")
}

