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
	fs.StringVar(&cfgPath, "config", "", "Pfad zur YAML-Konfigurationsdatei")

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
	default:
		return fmt.Errorf("unbekannter Befehl %q", rest[0])
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

	fs.StringVar(&text, "text", "", "Titel des Todos (Pflicht)")
	fs.Var(&checks, "check", "Ein Checklisten-Eintrag (kann mehrfach angegeben werden)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("Flag -text ist erforderlich")
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

	fmt.Fprintf(os.Stdout, "Todo angelegt: %s (ID: %s)\n", task.Text, task.ID)
	if len(task.Checklist) > 0 {
		fmt.Fprintln(os.Stdout, "Checkliste:")
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
		fmt.Fprintln(os.Stdout, "Keine Todos gefunden.")
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
	fs.StringVar(&id, "id", "", "ID der zu löschenden Todo (Pflicht)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("Flag -id ist erforderlich")
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

	fmt.Fprintf(os.Stdout, "Todo mit ID %s wurde gelöscht.\n", id)
	return nil
}

