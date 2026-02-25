package habitica

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTasksService_ListUserTasks(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks/user", r.URL.Path)
		require.Equal(t, "habits", r.URL.Query().Get("type"))

		resp := APIResponse[[]*Task]{
			Success: true,
			Data: []*Task{
				{
					ID:   "task-id-1",
					Text: "My Habit",
					Type: TaskTypeHabit,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	tasks, err := client.Tasks.ListUserTasks(context.Background(), TasksFilter{Type: "habits"})
	require.NoError(t, err)
	require.Len(t, tasks, 1)
	require.Equal(t, UUID("task-id-1"), tasks[0].ID)
	require.Equal(t, TaskTypeHabit, tasks[0].Type)
}

func TestTasksService_CreateTask(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks/user", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)

		var body TaskCreateRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "New Todo", body.Text)

		resp := APIResponse[Task]{
			Success: true,
			Data: Task{
				ID:   "task-id-2",
				Text: body.Text,
				Type: TaskTypeTodo,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	task, err := client.Tasks.CreateTask(context.Background(), &TaskCreateRequest{
		Text: "New Todo",
		Type: TaskTypeTodo,
	})
	require.NoError(t, err)
	require.Equal(t, UUID("task-id-2"), task.ID)
	require.Equal(t, "New Todo", task.Text)
}

func TestTasksService_CreateTodoWithChecklist(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks/user", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)

		var body TaskCreateRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.Equal(t, "My Todo", body.Text)
		require.Len(t, body.Checklist, 2)
		require.Equal(t, "Sub 1", body.Checklist[0].Text)
		require.False(t, body.Checklist[0].Completed)

		resp := APIResponse[Task]{
			Success: true,
			Data: Task{
				ID:        "task-id-3",
				Text:      body.Text,
				Type:      TaskTypeTodo,
				Checklist: body.Checklist,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	task, err := client.Tasks.CreateTodoWithChecklist(context.Background(), "My Todo", []string{"Sub 1", "Sub 2"})
	require.NoError(t, err)
	require.Equal(t, UUID("task-id-3"), task.ID)
	require.Len(t, task.Checklist, 2)
}

func TestTasksService_ScoreTask(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks/task-id-4/score/up", r.URL.Path)
		require.Equal(t, http.MethodPost, r.Method)

		resp := APIResponse[struct{}]{
			Success: true,
			Data:    struct{}{},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	err := client.Tasks.ScoreTask(context.Background(), UUID("task-id-4"), "up")
	require.NoError(t, err)
}

func TestTasksService_ScoreChecklistItem(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/tasks/task-id-5/checklist/item-id-1", r.URL.Path)
		require.Equal(t, http.MethodPut, r.Method)

		var body struct {
			Completed bool `json:"completed"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.True(t, body.Completed)

		resp := APIResponse[Task]{
			Success: true,
			Data: Task{
				ID:   "task-id-5",
				Text: "My Todo",
				Type: TaskTypeTodo,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	client, srv := newTestClient(t, handler)
	defer srv.Close()

	err := client.Tasks.UpdateChecklistItemCompleted(context.Background(), UUID("task-id-5"), UUID("item-id-1"), true)
	require.NoError(t, err)
}

