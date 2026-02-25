package habitica

// UserHistoryEntry describes a single entry of the user history
// as returned by /admin/user/:userId/history.
type UserHistoryEntry struct {
	Timestamp Timestamp   `json:"timestamp"`
	Data      map[string]any `json:"data"`
}

