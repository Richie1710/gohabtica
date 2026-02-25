package habitica

// GroupType describes the type of a group (party, guild, tavern).
type GroupType string

const (
	GroupTypeParty GroupType = "party"
	GroupTypeGuild GroupType = "guild"
	GroupTypeTavern GroupType = "tavern"
)

// Group represents a Habitica group.
type Group struct {
	ID          UUID            `json:"_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Type        GroupType       `json:"type"`
	LeaderID    UUID            `json:"leader"`
	MemberCount int             `json:"memberCount"`
	Privacy     string          `json:"privacy"`
}

