# Carpet Cleaner TUI Refactor Design

## Overview

Refactor `carpet-cleaner` from a comma-separated CLI input flow into a Bubble Tea TUI with selectable repos, supporting both delete and rename operations.

## Architecture

Single `tea.Model` with an internal state enum:

```
Loading → Select → Confirm → Executing → Done
```

- **Select** and **Confirm** render as labeled tabs at the top of the screen.
- Rename name entry uses a modal overlay.
- All state lives in one struct — no sub-models needed.

### Key packages

- `github.com/charmbracelet/bubbletea` — TUI loop and model
- `github.com/charmbracelet/bubbles/list` — scrollable, filterable repo list
- `github.com/charmouncelet/bubbles/textinput` — rename modal input
- `github.com/charmbracelet/bubbles/viewport` — execution results
- `github.com/charmbracelet/lipgloss` — styling and layout

### File structure

```
main.go            — entry point, env loading, API client functions
tui/
  model.go         — main model struct and Init/Update/View
  select.go        — select screen rendering and key handling
  confirm.go       — confirm screen rendering and key handling
  executing.go     — execution screen with live progress
  rename_modal.go  — modal overlay for entering new repo names
  styles.go        — lipgloss style definitions
```

## Screens

### Loading

Shown while fetching repos from GitHub API (paginated). Displays a spinner and "Fetching repositories..." text. Transitions to Select once all pages are fetched.

### Select Screen

- Scrollable list of all fetched repos
- Each item renders as a custom `list.Item` whose `Title()` returns `[ ] owner/repo` (or `[x]` when selected) and `Description()` returns the repo description and private/public badge. The checkbox state is tracked in the model, not in the delegate.
- `↑`/`↓` to navigate, `Space`/`Enter` to toggle selection
- `/` to activate filter (built into `bubbles/list`)
- `d` sets operation to Delete, `r` sets to Rename
- Status bar: "N selected, operation: delete" or "N selected, operation: rename"
- `Tab` or `Enter` (when items are selected) moves to Confirm

### Rename Modal

When operation is Rename and user advances to Confirm, a modal dialog pops up sequentially for each selected repo:

- Shows current repo name
- `textinput` field pre-filled with the current name
- Enter confirms the new name, Escape skips/cancels that repo
- If all renames are cancelled, returns to Select

### Confirm Screen

- Tab bar: `Select | Confirm` (Confirm highlighted)
- Summary of pending operations:
  - "Delete: repo1, repo2"
  - "Rename: repo3 → new-name"
- Warning banner for delete operations
- `Enter` to execute, `Esc` to go back to Select

### Executing Screen

- Runs operations sequentially
- Live progress per operation:
  - Spinner + "Deleting repo1..." → ✓ or ✗ with error message
  - Spinner + "Renaming repo2 → new-name..." → ✓ or ✗ with error message
- Scrollable viewport so output doesn't disappear
- After all operations complete: summary — "3 deleted, 1 renamed, 1 failed. Press q to quit."

## Key Bindings Summary

| Key       | Select Screen        | Confirm Screen | Executing Screen |
|-----------|----------------------|----------------|------------------|
| ↑/↓       | Navigate list        | —              | —                |
| Space     | Toggle selection     | —              | —                |
| /         | Filter               | —              | —                |
| d         | Set operation: delete| —              | —                |
| r         | Set operation: rename| —              | —                |
| Enter     | Proceed to confirm  | Execute        | —                |
| Esc       | Quit (if none selected) | Back to select| —             |
| q         | —                    | —              | Quit             |

## Bug Fixes (resolved by refactor)

These existing bugs are resolved as natural consequences of the TUI refactor:

- **#1 RenameRepo always returns error** — Fix: add success status code check (`200 OK` → return nil)
- **#2 RenameRepo never called** — Fix: wired into TUI via rename operation
- **#4 No validation of repo names** — Fix: repos come from fetched list, no manual name entry
- **#5 Trailing comma empty entries** — Fix: no comma-separated input at all

Remaining bugs not addressed by this refactor (separate work):
- **#3 Content-Type header inconsistency** — cosmetic, no functional impact
- **#6 No rate-limit handling** — Bubble Tea's async model makes this easier, but not in scope
- **#7 godotenv error silently ignored** — not a TUI concern, but worth a quick fix

## Error Handling

- API errors during fetch: show error in Loading screen with "Press q to quit"
- API errors during execution: display inline per operation with ✗ marker, continue to next operation
- No network: caught at fetch stage, no TUI rendered