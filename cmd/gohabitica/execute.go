package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/danielrichardt/gohabitica/habitica"
	"github.com/danielrichardt/gohabitica/internal/config"
)

// stringSliceFlag allows passing a flag multiple times.
type stringSliceFlag []string

func (s *stringSliceFlag) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSliceFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
}

// execute is the entry point for the CLI logic.
// Without arguments, it runs a simple smoke test (GET /user).
// With subcommand "todo" it creates todos with checklists.
// With subcommand "todos" it lists existing todos.
// With the "-config" flag you can specify an explicit YAML configuration file.
func execute(args []string) error {
	fs := flag.NewFlagSet("gohabitica", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var cfgPath string
	fs.StringVar(&cfgPath, "config", "", "Path to a YAML configuration file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	rest := fs.Args()

	if len(rest) == 0 {
		return runSmokeTest(cfgPath)
	}

	switch rest[0] {
	case "todo":
		return runTodo(cfgPath, rest[1:])
	case "todos":
		return runTodosList(cfgPath, rest[1:])
	case "todo-delete":
		return runTodoDelete(cfgPath, rest[1:])
	case "todo-complete":
		return runTodoComplete(cfgPath, rest[1:])
	case "todo-check":
		return runTodoCheck(cfgPath, rest[1:])
	default:
		return fmt.Errorf("unknown command %q", rest[0])
	}
}

// runSmokeTest runs the original /user sanity check.
func runSmokeTest(cfgPath string) error {
	opts := config.Options{ConfigPath: cfgPath}
	cfg, err := config.Load(opts)
	if err != nil {
		return err
	}

	client, err := habitica.NewClient(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u, err := client.User.GetCurrent(ctx)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stdout, "Eingeloggt als: %s\n", u.Profile.Name)
	return err
}

// runTodo creates a new todo (TaskTypeTodo) with an optional checklist.
//
// Example usage:
//   gohabitica todo -text "Shopping" -check "Milk" -check "Bread"
//   gohabitica -config config/local.yaml todo -text "Shopping" -check "Milk"
func runTodo(cfgPath string, args []string) error {
	fs := flag.NewFlagSet("todo", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var text string
	var checks stringSliceFlag

	fs.StringVar(&text, "text", "", "Title of the todo (required)")
	fs.Var(&checks, "check", "A checklist entry (can be specified multiple times)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("flag -text is required")
	}

	opts := config.Options{ConfigPath: cfgPath}
	cfg, err := config.Load(opts)
	if err != nil {
		return err
	}

	client, err := habitica.NewClient(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	task, err := client.Tasks.CreateTodoWithChecklist(ctx, text, checks)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Todo created: %s (ID: %s)\n", task.Text, task.ID)
	if len(task.Checklist) > 0 {
		fmt.Fprintln(os.Stdout, "Checklist:")
		for _, item := range task.Checklist {
			fmt.Fprintf(os.Stdout, "  - %s\n", item.Text)
		}
	}

	return nil
}

// runTodosList lists the user's todos together with their checklist entries.
//
// Example usage:
//   gohabitica todos
//   gohabitica -config config/local.yaml todos
func runTodosList(cfgPath string, args []string) error {
	// currently no flags; can be extended later if needed
	_ = args

	opts := config.Options{ConfigPath: cfgPath}
	cfg, err := config.Load(opts)
	if err != nil {
		return err
	}

	client, err := habitica.NewClient(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tasks, err := client.Tasks.ListUserTasks(ctx, habitica.TasksFilter{Type: "todos"})
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		fmt.Fprintln(os.Stdout, "No todos found.")
		return nil
	}

	for _, t := range tasks {
		status := " "
		if t.Completed {
			status = "x"
		}
		fmt.Fprintf(os.Stdout, "[%s] %s (ID: %s)\n", status, t.Text, t.ID)
		if len(t.Checklist) > 0 {
			for _, item := range t.Checklist {
				subStatus := " "
				if item.Completed {
					subStatus = "x"
				}
				fmt.Fprintf(os.Stdout, "  - [%s] %s\n", subStatus, item.Text)
			}
		}
	}

	return nil
}

// runTodoDelete deletes a todo by its ID.
//
// Example usage:
//   gohabitica todo-delete -id "37ceed6f-0772-43bb-a177-39d3074f75b7"
//   gohabitica -config config/local.yaml todo-delete -id "..."
func runTodoDelete(cfgPath string, args []string) error {
	fs := flag.NewFlagSet("todo-delete", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var id string
	fs.StringVar(&id, "id", "", "ID of the todo to delete (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("flag -id is required")
	}

	opts := config.Options{ConfigPath: cfgPath}
	cfg, err := config.Load(opts)
	if err != nil {
		return err
	}

	client, err := habitica.NewClient(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Tasks.DeleteTask(ctx, habitica.UUID(id)); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Todo with ID %s has been deleted.\n", id)
	return nil
}

// runTodoComplete marks a todo as completed by scoring it "up".
//
// Example usage:
//   gohabitica todo-complete -id "37ceed6f-0772-43bb-a177-39d3074f75b7"
//   gohabitica -config config/local.yaml todo-complete -id "..."
func runTodoComplete(cfgPath string, args []string) error {
	fs := flag.NewFlagSet("todo-complete", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var id string
	fs.StringVar(&id, "id", "", "ID of the todo to complete (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("flag -id is required")
	}

	opts := config.Options{ConfigPath: cfgPath}
	cfg, err := config.Load(opts)
	if err != nil {
		return err
	}

	client, err := habitica.NewClient(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Tasks.ScoreTask(ctx, habitica.UUID(id), "up"); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Todo with ID %s has been completed.\n", id)
	return nil
}

// runTodoCheck toggles a checklist item of a todo by its 1-based index.
//
// Example usage:
//   gohabitica todo-check -id "37ceed6f-0772-43bb-a177-39d3074f75b7" -index 1
//   gohabitica -config config/local.yaml todo-check -id "..." -index 2
func runTodoCheck(cfgPath string, args []string) error {
	fs := flag.NewFlagSet("todo-check", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var (
		id    string
		index int
	)

	fs.StringVar(&id, "id", "", "ID of the todo that owns the checklist item (required)")
	fs.IntVar(&index, "index", 0, "1-based index of the checklist item to toggle (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("flag -id is required")
	}
	if index <= 0 {
		return fmt.Errorf("flag -index must be greater than zero")
	}

	opts := config.Options{ConfigPath: cfgPath}
	cfg, err := config.Load(opts)
	if err != nil {
		return err
	}

	client, err := habitica.NewClient(cfg)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	task, err := client.Tasks.GetTask(ctx, habitica.UUID(id))
	if err != nil {
		return err
	}
	if task.Type != habitica.TaskTypeTodo {
		return fmt.Errorf("task %s is not a todo", id)
	}
	if len(task.Checklist) == 0 {
		return fmt.Errorf("todo %s has no checklist items", id)
	}
	if index > len(task.Checklist) {
		return fmt.Errorf("flag -index is out of range: todo %s only has %d checklist items", id, len(task.Checklist))
	}

	item := task.Checklist[index-1]
	if err := client.Tasks.UpdateChecklistItemCompleted(ctx, task.ID, item.ID, !item.Completed); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "Checklist item #%d (%q) on todo %s has been toggled.\n", index, item.Text, task.ID)
	return nil
}

