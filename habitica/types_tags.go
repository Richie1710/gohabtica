package habitica

// Tag represents a tag that can be used to group tasks.
type Tag struct {
	ID   UUID   `json:"id"`
	Name string `json:"name"`
}

