package habitica

import (
	"context"
	"fmt"
	"net/url"
)

// TasksService groups endpoints related to user tasks.
type TasksService struct {
	client *Client
}

// TasksFilter restricts the result set of ListUserTasks.
type TasksFilter struct {
	Type string // habits, dailys, todos, rewards, completedTodos
}

// ListUserTasks retrieves the tasks of the current user (GET /tasks/user).
func (s *TasksService) ListUserTasks(ctx context.Context, filter TasksFilter) ([]*Task, error) {
	q := url.Values{}
	if filter.Type != "" {
		q.Set("type", filter.Type)
	}
	var tasks []*Task
	if err := s.client.doRequest(ctx, "GET", "/tasks/user", q, nil, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// GetTask retrieves a single task by its ID (GET /tasks/:id).
func (s *TasksService) GetTask(ctx context.Context, id UUID) (*Task, error) {
	var t Task
	if err := s.client.doRequest(ctx, "GET", fmt.Sprintf("/tasks/%s", id), nil, nil, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// CreateTask creates a new task (POST /tasks/user).
func (s *TasksService) CreateTask(ctx context.Context, in *TaskCreateRequest) (*Task, error) {
	var t Task
	if err := s.client.doRequest(ctx, "POST", "/tasks/user", nil, in, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// CreateTodoWithChecklist is a convenience wrapper to create a todo with a checklist
// from plain strings.
func (s *TasksService) CreateTodoWithChecklist(ctx context.Context, text string, checklist []string) (*Task, error) {
	items := make([]ChecklistItem, 0, len(checklist))
	for _, c := range checklist {
		if c == "" {
			continue
		}
		items = append(items, ChecklistItem{
			Text:      c,
			Completed: false,
		})
	}

	req := &TaskCreateRequest{
		Text:      text,
		Type:      TaskTypeTodo,
		Priority:  1,
		Attribute: "str",
		Checklist: items,
	}
	return s.CreateTask(ctx, req)
}

// UpdateTask updates an existing task (PUT /tasks/:id).
func (s *TasksService) UpdateTask(ctx context.Context, id UUID, in *TaskUpdateRequest) (*Task, error) {
	var t Task
	if err := s.client.doRequest(ctx, "PUT", fmt.Sprintf("/tasks/%s", id), nil, in, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

// DeleteTask deletes a task (DELETE /tasks/:id).
func (s *TasksService) DeleteTask(ctx context.Context, id UUID) error {
	return s.client.doRequest(ctx, "DELETE", fmt.Sprintf("/tasks/%s", id), nil, nil, nil)
}

