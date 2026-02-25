package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupsService_GetGroup(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/groups/group-1", r.URL.Path)

		resp := APIResponse[Group]{
			Success: true,
			Data: Group{
				ID:   "group-1",
				Name: "My Party",
				Type: GroupTypeParty,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	g, err := client.Groups.GetGroup(context.Background(), "group-1")
	require.NoError(t, err)
	require.Equal(t, UUID("group-1"), g.ID)
	require.Equal(t, GroupTypeParty, g.Type)
}

