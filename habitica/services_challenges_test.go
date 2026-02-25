package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChallengesService_GetChallenge(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/challenges/chal-1", r.URL.Path)

		resp := APIResponse[Challenge]{
			Success: true,
			Data: Challenge{
				ID:   "chal-1",
				Name: "My Challenge",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	ch, err := client.Challenges.GetChallenge(context.Background(), "chal-1")
	require.NoError(t, err)
	require.Equal(t, UUID("chal-1"), ch.ID)
	require.Equal(t, "My Challenge", ch.Name)
}

