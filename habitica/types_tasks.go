package habitica

// TaskType describes the type of a task.
type TaskType string

const (
	TaskTypeHabit  TaskType = "habit"
	TaskTypeDaily  TaskType = "daily"
	TaskTypeTodo   TaskType = "todo"
	TaskTypeReward TaskType = "reward"
)

// Task represents a Habitica task.
type Task struct {
	ID          UUID            `json:"_id"`
	UserID      UUID            `json:"userId"`
	Text        string          `json:"text"`
	Notes       string          `json:"notes"`
	Type        TaskType        `json:"type"`
	Value       float64         `json:"value"`
	Priority    float64         `json:"priority"`
	CreatedAt   Timestamp       `json:"createdAt"`
	UpdatedAt   Timestamp       `json:"updatedAt"`
	Tags        []UUID          `json:"tags"`
	Completed   bool            `json:"completed"`
	Checklist   []ChecklistItem `json:"checklist"`
	Attribute   string          `json:"attribute"` // str, int, con, per
}

// ChecklistItem is an entry of a todo or daily checklist.
type ChecklistItem struct {
	ID        UUID   `json:"id,omitempty"`
	Text      string `json:"text"`
	Completed bool   `json:"completed,omitempty"`
}

// TaskCreateRequest describes the fields used to create a new task.
type TaskCreateRequest struct {
	Text      string          `json:"text"`
	Notes     string          `json:"notes,omitempty"`
	Type      TaskType        `json:"type"`
	Priority  float64         `json:"priority,omitempty"`
	Tags      []UUID          `json:"tags,omitempty"`
	Checklist []ChecklistItem `json:"checklist,omitempty"`
	Attribute string          `json:"attribute,omitempty"`
}

// TaskUpdateRequest describes the fields used to update an existing task.
type TaskUpdateRequest struct {
	Text      *string          `json:"text,omitempty"`
	Notes     *string          `json:"notes,omitempty"`
	Priority  *float64         `json:"priority,omitempty"`
	Tags      *[]UUID          `json:"tags,omitempty"`
	Checklist *[]ChecklistItem `json:"checklist,omitempty"`
	Attribute *string          `json:"attribute,omitempty"`
}

