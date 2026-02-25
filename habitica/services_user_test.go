package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserService_GetCurrent(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/user", r.URL.Path)
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "user-id", r.Header.Get("x-api-user"))
		require.Equal(t, "api-token", r.Header.Get("x-api-key"))

		resp := APIResponse[User]{
			Success: true,
			Data: User{
				ID: "user-id",
				Profile: UserProfile{
					Name: "Test User",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	u, err := client.User.GetCurrent(context.Background())
	require.NoError(t, err)
	require.NotNil(t, u)
	require.Equal(t, UUID("user-id"), u.ID)
	require.Equal(t, "Test User", u.Profile.Name)
}

