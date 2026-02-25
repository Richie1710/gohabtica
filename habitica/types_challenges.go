package habitica

// Challenge represents a Habitica challenge.
type Challenge struct {
	ID          UUID        `json:"_id"`
	Name        string      `json:"name"`
	ShortName   string      `json:"shortName"`
	Description string      `json:"description"`
	LeaderID    UUID        `json:"leader"`
	GroupID     UUID        `json:"group"`
	MemberCount int         `json:"memberCount"`
	Prize       int         `json:"prize"`
	Active      bool        `json:"isActive"`
}

