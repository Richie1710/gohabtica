package habitica

// User represents a Habitica user (reduced to the fields typically needed by clients,
// structured according to the official API documentation).
type User struct {
	ID          UUID             `json:"_id"`
	Auth        UserAuth         `json:"auth"`
	Profile     UserProfile      `json:"profile"`
	Stats       UserStats        `json:"stats"`
	Preferences UserPreferences  `json:"preferences"`
	Items       UserItems        `json:"items"`
	Flags       map[string]any   `json:"flags"`
	Notifications []UserNotification `json:"notifications"`
}

type UserAuth struct {
	Local struct {
		Email string `json:"email"`
	} `json:"local"`
}

type UserProfile struct {
	Name     string `json:"name"`
	Blurb    string `json:"blurb"`
	ImageURL string `json:"imageUrl"`
}

type UserStats struct {
	Class           string   `json:"class"`
	HP              float64  `json:"hp"`
	MP              float64  `json:"mp"`
	Exp             float64  `json:"exp"`
	GP              float64  `json:"gp"`
	Lvl             int      `json:"lvl"`
	Points          int      `json:"points"`
	MaxHP           float64  `json:"maxHealth"`
	MaxMP           float64  `json:"maxMP"`
	Buffs           map[string]any `json:"buffs"`
	ToNextLevel     float64  `json:"toNextLevel"`
}

type UserPreferences struct {
	DayStart      int    `json:"dayStart"`
	Language      string `json:"language"`
	Background    string `json:"background"`
	TimezoneOffset float64 `json:"timezoneOffset"`
}

type UserItems struct {
	Gold           float64          `json:"gold"`
	Gems           int              `json:"gems"`
	Equipment      map[string]any   `json:"gear"`
	Pets           map[string]int   `json:"pets"`
	Mounts         map[string]bool  `json:"mounts"`
	CurrentMount   string           `json:"currentMount"`
	CurrentPet     string           `json:"currentPet"`
}

type UserNotification struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Seen      bool           `json:"seen"`
	Data      map[string]any `json:"data"`
	Timestamp Timestamp      `json:"timestamp"`
}

