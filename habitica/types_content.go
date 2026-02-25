package habitica

// Content represents the static content endpoint (/content).
// The structure is intentionally simplified; in most clients, it is
// sufficient to read only specific sub-sections as needed.
type Content struct {
	Items    map[string]any `json:"items"`
	Pets     map[string]any `json:"pets"`
	Mounts   map[string]any `json:"mounts"`
	Quests   map[string]any `json:"quests"`
	Backgrounds map[string]any `json:"backgrounds"`
}

