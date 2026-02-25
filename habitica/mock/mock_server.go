package mock

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/danielrichardt/gohabitica/habitica"
	"github.com/danielrichardt/gohabitica/internal/config"
)

// Server wraps an httptest.Server together with a fully configured Habitica client.
type Server struct {
	HTTP   *httptest.Server
	Client *habitica.Client
}

// NewServer creates a new mock server with the given handler.
func NewServer(handler http.Handler) (*Server, error) {
	srv := httptest.NewServer(handler)

	cfg := &config.Config{
		BaseURL:  srv.URL,
		UserID:   "test-user",
		APIToken: "test-token",
	}

	client, err := habitica.NewClient(cfg, habitica.WithHTTPClient(&http.Client{
		Timeout: 5 * time.Second,
	}))
	if err != nil {
		srv.Close()
		return nil, err
	}

	return &Server{
		HTTP:   srv,
		Client: client,
	}, nil
}

// Close shuts down the underlying HTTP server.
func (s *Server) Close() {
	if s.HTTP != nil {
		s.HTTP.Close()
	}
}

