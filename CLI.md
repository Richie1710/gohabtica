## Gohabitica CLI – Reference (Human + AI Friendly)

This document describes how to use the `gohabitica` CLI to interact with the Habitica API.
It is written to be easy to parse for both humans and language models.

---

### 1. Overview

- **Purpose**: The `gohabitica` CLI connects to the Habitica API and allows you to:
  - Test connectivity (`GET /user`)
  - Create todos with optional checklists
  - List existing todos and their checklist items
  - Delete todos
  - Mark todos as completed
  - Toggle checklist items on a todo

- **Binary name**: `gohabitica`
- **Default behavior (no subcommand)**: Runs a smoke test (`GET /user`) and prints the logged-in user.

---

### 2. Configuration

The CLI needs Habitica credentials: **User ID** and **API Token**.

Configuration can be provided in **two ways**:

#### 2.1 Environment variables (preferred)

Set these environment variables:

- `HABITICA_USER_ID` – your Habitica User ID
- `HABITICA_API_TOKEN` – your Habitica API token

If both variables are set, the CLI will use them directly and **no config file is required**.

Example (bash/zsh):

```bash
export HABITICA_USER_ID="your-user-id-here"
export HABITICA_API_TOKEN="your-api-token-here"
gohabitica
```

#### 2.2 YAML configuration file

If the environment variables are **not** set, the CLI looks for a YAML file.

- **Default location** (derived from `os.UserConfigDir()` in Go):

  `<user-config-dir>/gohabitica/config.yaml`

  Typical values:
  - Linux: `$XDG_CONFIG_HOME/gohabitica/config.yaml` or `$HOME/.config/gohabitica/config.yaml`
  - macOS: `$HOME/Library/Application Support/gohabitica/config.yaml`
  - Windows: `%AppData%\gohabitica\config.yaml`

- **Override location** via global flag `-config`:

  ```bash
  gohabitica -config /path/to/config.yaml <subcommand> ...
  ```

**YAML structure:**

```yaml
base_url: https://habitica.com/api/v3   # optional, defaults to Habitica production API
user_id: "your-user-id-here"           # required
api_token: "your-api-token-here"       # required
```

Rules:

- If `user_id` or `api_token` is missing, the CLI returns an error about missing credentials.
- If `base_url` is omitted, it defaults to `https://habitica.com/api/v3`.

---

### 3. Global flags

Global flags must appear **before** the subcommand.

- **Flag**: `-config <path>`
  - **Type**: string (file path)
  - **Purpose**: Path to a YAML configuration file (see section 2.2).
  - **Required**: No
  - **Example**:
    ```bash
    gohabitica -config config/local.yaml todos
    ```

No other global flags are currently supported.

---

### 4. Commands

#### 4.1 No subcommand (smoke test)

- **Synopsis**:
  ```bash
  gohabitica [ -config <path> ]
  ```
- **Description**:
  Performs a simple smoke test by calling `GET /user` and prints the logged-in user.

- **Output example**:
  ```text
  Eingeloggt als: Alice Example
  ```

- **Errors / exit code**:
  - Returns a non-zero exit code if:
    - Credentials are missing or invalid
    - The Habitica API is not reachable
    - Any unexpected error occurs

---

#### 4.2 `todo` – create a new todo with optional checklist

- **Synopsis**:
  ```bash
  gohabitica [ -config <path> ] todo -text "<title>" [ -check "<item1>" ... ]
  ```

- **Description**:
  Creates a new Habitica **todo** (`TaskTypeTodo`) with the given title and an optional checklist.

- **Flags (command-level)**:
  - `-text <string>`
    - **Type**: string
    - **Required**: Yes
    - **Meaning**: Title of the todo.
  - `-check <string>`
    - **Type**: string, **repeatable**
    - **Required**: No
    - **Meaning**: A checklist entry. Can be specified multiple times to create multiple checklist items.

- **Examples**:
  ```bash
  # Simple todo without checklist
  gohabitica todo -text "Buy groceries"

  # Todo with multiple checklist items
  gohabitica todo -text "Buy groceries" \
    -check "Milk" \
    -check "Bread" \
    -check "Eggs"

  # Using explicit config file
  gohabitica -config config/local.yaml todo -text "Plan week" -check "Review calendar"
  ```

- **Output example**:
  ```text
  Todo created: Buy groceries (ID: 37ceed6f-0772-43bb-a177-39d3074f75b7)
  Checklist:
    - Milk
    - Bread
    - Eggs
  ```

- **Errors / exit code**:
  - `-text` empty or missing → error: `flag -text is required`
  - API / network errors → non-zero exit code.

---

#### 4.3 `todos` – list all todos and their checklists

- **Synopsis**:
  ```bash
  gohabitica [ -config <path> ] todos
  ```

- **Description**:
  Lists all Habitica tasks of type `"todos"` for the current user, including checklist items and completion status.

- **Flags (command-level)**:
  - Currently none.

- **Behavior**:
  - Fetches all user tasks filtered by type `todos`.
  - For each todo, prints:
    - Completion status `[ ]` or `[x]`
    - Text (title)
    - Internal Habitica ID
  - For each checklist item, prints:
    - Indented checklist lines with `[ ]` or `[x]` and the text.

- **Examples**:
  ```bash
  gohabitica todos
  gohabitica -config config/local.yaml todos
  ```

- **Output examples**:
  ```text
  No todos found.
  ```

  or:

  ```text
  [ ] Buy groceries (ID: 37ceed6f-0772-43bb-a177-39d3074f75b7)
    - [ ] Milk
    - [x] Bread
    - [ ] Eggs
  [x] Plan week (ID: 0a123456-789b-4cde-f012-3456789abcde)
  ```

---

#### 4.4 `todo-delete` – delete a todo by ID

- **Synopsis**:
  ```bash
  gohabitica [ -config <path> ] todo-delete -id "<todo-id>"
  ```

- **Description**:
  Deletes an existing todo identified by its Habitica task ID.

- **Flags (command-level)**:
  - `-id <string>`
    - **Required**: Yes
    - **Meaning**: ID of the todo to delete.

- **Examples**:
  ```bash
  gohabitica todo-delete -id "37ceed6f-0772-43bb-a177-39d3074f75b7"
  gohabitica -config config/local.yaml todo-delete -id "0a123456-789b-4cde-f012-3456789abcde"
  ```

- **Output example**:
  ```text
  Todo with ID 37ceed6f-0772-43bb-a177-39d3074f75b7 has been deleted.
  ```

- **Errors / exit code**:
  - `-id` empty or missing → error: `flag -id is required`
  - ID does not exist or cannot be deleted → non-zero exit code.

---

#### 4.5 `todo-complete` – mark a todo as completed

- **Synopsis**:
  ```bash
  gohabitica [ -config <path> ] todo-complete -id "<todo-id>"
  ```

- **Description**:
  Marks an existing todo as completed by scoring it `"up"` via the Habitica API.

- **Flags (command-level)**:
  - `-id <string>`
    - **Required**: Yes
    - **Meaning**: ID of the todo to complete.

- **Examples**:
  ```bash
  gohabitica todo-complete -id "37ceed6f-0772-43bb-a177-39d3074f75b7"
  gohabitica -config config/local.yaml todo-complete -id "0a123456-789b-4cde-f012-3456789abcde"
  ```

- **Output example**:
  ```text
  Todo with ID 37ceed6f-0772-43bb-a177-39d3074f75b7 has been completed.
  ```

- **Errors / exit code**:
  - `-id` empty or missing → error: `flag -id is required`
  - ID invalid or API error → non-zero exit code.

---

#### 4.6 `todo-check` – toggle a checklist item on a todo

- **Synopsis**:
  ```bash
  gohabitica [ -config <path> ] todo-check -id "<todo-id>" -index <n>
  ```

- **Description**:
  Toggles the completion state of a checklist item on a todo.
  The checklist item is selected by a **1-based** index.

- **Flags (command-level)**:
  - `-id <string>`
    - **Required**: Yes
    - **Meaning**: ID of the todo owning the checklist.
  - `-index <int>`
    - **Required**: Yes
    - **Meaning**: 1-based index of the checklist item to toggle.

- **Behavior**:
  1. Fetches the task by ID.
  2. Verifies that the task is a todo.
  3. Verifies that there is a checklist.
  4. Verifies that `index` is within range.
  5. Flips the `completed` flag of the selected checklist item.

- **Examples**:
  ```bash
  # Toggle first checklist item
  gohabitica todo-check -id "37ceed6f-0772-43bb-a177-39d3074f75b7" -index 1

  # Toggle second checklist item using explicit config
  gohabitica -config config/local.yaml todo-check -id "37ceed6f-0772-43bb-a177-39d3074f75b7" -index 2
  ```

- **Output example**:
  ```text
  Checklist item #1 ("Milk") on todo 37ceed6f-0772-43bb-a177-39d3074f75b7 has been toggled.
  ```

- **Errors / exit code**:
  - `-id` empty → error: `flag -id is required`
  - `-index` <= 0 → error: `flag -index must be greater than zero`
  - Task not found or not a todo → error
  - Checklist empty or index out of range → error with informative message.

---

### 5. Machine-readable command summary

This section summarizes commands and flags in a structure that is easy for tools and language models to parse.

- **Global**
  - **Flag**: `-config <string>` – optional, path to YAML config file.

- **Command**: (no subcommand)
  - **Purpose**: smoke test (`GET /user`)
  - **Args**: none (only global flags)

- **Command**: `todo`
  - **Purpose**: create a new todo with optional checklist.
  - **Flags**:
    - `-text <string>` – required
    - `-check <string>` – optional, repeatable

- **Command**: `todos`
  - **Purpose**: list todos and their checklist items.
  - **Flags**: none

- **Command**: `todo-delete`
  - **Purpose**: delete a todo by ID.
  - **Flags**:
    - `-id <string>` – required

- **Command**: `todo-complete`
  - **Purpose**: mark a todo as completed.
  - **Flags**:
    - `-id <string>` – required

- **Command**: `todo-check`
  - **Purpose**: toggle a checklist item.
  - **Flags**:
    - `-id <string>` – required
    - `-index <int>` – required, 1-based

