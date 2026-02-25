package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/danielrichardt/gohabitica/internal/config"
	"github.com/stretchr/testify/require"
)

// newTestClient creates a client wired up against a mocked HTTP server.
func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()

	srv := httptest.NewServer(handler)

	cfg := &config.Config{
		BaseURL:  srv.URL,
		UserID:   "user-id",
		APIToken: "api-token",
	}

	c, err := NewClient(cfg)
	require.NoError(t, err)

	return c, srv
}

func TestClient_DoRequest_Success(t *testing.T) {
	type payload struct {
		Foo string `json:"foo"`
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(APIResponse[payload]{
			Success: true,
			Data:    payload{Foo: "bar"},
		})
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	var out payload
	err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil, &out)
	require.NoError(t, err)
	require.Equal(t, "bar", out.Foo)
}

func TestClient_DoRequest_Error(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(APIResponse[map[string]any]{
			Success: false,
			Error:   "NotAuthorized",
			Message: "Missing authentication headers.",
		})
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	err := client.doRequest(context.Background(), http.MethodGet, "/test", nil, nil, nil)
	require.Error(t, err)

	var apiErr *APIError
	require.ErrorAs(t, err, &apiErr)
	require.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
	require.Equal(t, "NotAuthorized", apiErr.Code)
	require.True(t, IsUnauthorized(err))
}

