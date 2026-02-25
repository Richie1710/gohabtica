# gohabtica

An idiomatic Go client and CLI for the [Habitica](https://habitica.com/) API.

## Features

- Typed Go client for the Habitica v3 API.
- Support for core domains: users, tasks, groups, challenges, content, tags, shops, webhooks, admin.
- Convenient helpers for todos including checklists.
- Simple CLI to experiment with your Habitica account.

## Installation

```bash
go get github.com/danielrichardt/gohabitica
```

Prebuilt CLI binaries for common platforms are published via GitHub Actions:

- On every push to `main` a snapshot build runs and uploads artifacts for Linux, macOS and Windows.
- On every tag starting with `v` (for example `v1.2.3`) a GitHub Release is created with uploaded archives.

You can download the latest binaries from the **Releases** page of the repository.

## Configuration

The client and CLI read credentials from either:

1. Environment variables:
   - `HABITICA_USER_ID`
   - `HABITICA_API_TOKEN`
2. Or a YAML config file.

For local development, you can use `config/local.yaml` (ignored by git):

```yaml
base_url: https://habitica.com/api/v3
user_id: "YOUR_HABITICA_USER_ID"
api_token: "YOUR_HABITICA_API_TOKEN"
```

Then pass it to the CLI via:

```bash
./gohabitica -config config/local.yaml ...
```

## CLI Usage

Build the CLI:

```bash
go build ./cmd/gohabitica
```

### Smoke test

Print the currently authenticated user:

```bash
./gohabitica -config config/local.yaml
```

### List todos with checklists

```bash
./gohabitica -config config/local.yaml todos
```

### Create a todo with checklist

```bash
./gohabitica -config config/local.yaml todo \
  -text "Clean the kitchen" \
  -check "Wipe counters" \
  -check "Sweep floor"
```

### Delete a todo

```bash
./gohabitica -config config/local.yaml todo-delete \
  -id "TASK_ID_FROM_TODOS_COMMAND"
```

## Library Usage (Example)

```go
ctx := context.Background()

cfg, err := config.Load(config.Options{})
if err != nil {
    log.Fatal(err)
}

client, err := habitica.NewClient(cfg)
if err != nil {
    log.Fatal(err)
}

tasks, err := client.Tasks.ListUserTasks(ctx, habitica.TasksFilter{Type: "todos"})
if err != nil {
    log.Fatal(err)
}

for _, t := range tasks {
    fmt.Println("Todo:", t.Text)
}
```
