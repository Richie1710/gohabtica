package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdminService_GetUserHistory(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/admin/user/test-user/history", r.URL.Path)

		resp := APIResponse[[]*UserHistoryEntry]{
			Success: true,
			Data: []*UserHistoryEntry{
				{Timestamp: "2024-01-01T00:00:00.000Z"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	hist, err := client.Admin.GetUserHistory(context.Background(), "test-user")
	require.NoError(t, err)
	require.Len(t, hist, 1)
	require.Equal(t, Timestamp("2024-01-01T00:00:00.000Z"), hist[0].Timestamp)
}

