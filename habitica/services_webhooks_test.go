package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWebhooksService_ListWebhooks(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/user/webhook", r.URL.Path)

		resp := APIResponse[[]*Webhook]{
			Success: true,
			Data: []*Webhook{
				{ID: "hook-1", URL: "https://example.com/hook"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	hooks, err := client.Webhooks.ListWebhooks(context.Background())
	require.NoError(t, err)
	require.Len(t, hooks, 1)
	require.Equal(t, UUID("hook-1"), hooks[0].ID)
}

