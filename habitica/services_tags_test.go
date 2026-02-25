package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTagsService_ListTags(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tags", r.URL.Path)

		resp := APIResponse[[]*Tag]{
			Success: true,
			Data: []*Tag{
				{ID: "tag-1", Name: "Work"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	tags, err := client.Tags.ListTags(context.Background())
	require.NoError(t, err)
	require.Len(t, tags, 1)
	require.Equal(t, UUID("tag-1"), tags[0].ID)
}

