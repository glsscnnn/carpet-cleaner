# Known Bugs

## 1. `RenameRepo` always returns an error

**File:** `main.go:53`

The `RenameRepo` function unconditionally returns an error on line 53:

```go
body, _ := io.ReadAll(resp.Body)
return fmt.Errorf("failed to update repo. Status: %d, Response: %s", resp.StatusCode, string(body))
```

This line runs regardless of the HTTP response status code. Even when the GitHub API returns a `200 OK` (successful rename), the function still returns a non-nil error with a "failed to update repo" message. It needs a success-path check (similar to `DeleteRepo`, which checks for `204 No Content` before returning `nil`).

**Fix:** Check for a successful status code before returning the error. A successful rename returns `200 OK`, so the function should return `nil` in that case:

```go
if resp.StatusCode == http.StatusOK {
    return nil
}
```

## 2. `RenameRepo` is defined but never called

**File:** `main.go:24-54`

The `RenameRepo` function exists but is never referenced in `main()` or anywhere else. The `main()` function only supports the delete operation — there is no prompt, menu, or flag to trigger a rename.

**Fix:** Wire `RenameRepo` into the CLI flow, e.g., by adding an operation selection step or a command-line flag to choose between delete and rename.

## 3. `RenameRepo` does not set `Content-Type` header correctly

**File:** `main.go:43`

While `Content-Type: application/json` is set on line 43, this is actually correct. However, note that `DeleteRepo` does **not** set `Content-Type` because it sends no request body. The inconsistency is fine functionally, but worth noting for consistency.

## 4. No validation that user-entered repo names actually exist

**File:** `main.go:169-194`

The user enters comma-separated repo names to delete, but there is no check that these names correspond to repos that were actually fetched. If a user typos a repo name, `DeleteRepo` will still attempt the API call and fail with a GitHub error, but the user won't get a helpful message like "repo not found in your account."

**Fix:** Validate each entered name against `allRepos` before attempting deletion, and warn the user about unrecognized names.

## 5. Trailing comma in input creates empty repo name entries

**File:** `main.go:169`

If the user enters a trailing comma (e.g., `repo1, repo2,`), `strings.Split` will produce an empty string as the last element. The code on line 186-188 skips empty strings with `if repo == "" { continue }`, so this doesn't cause a crash, but it does silently skip the entry without informing the user.

**Fix:** Trim trailing commas or warn the user about empty entries.

## 6. No pagination limit or rate-limit handling

**File:** `main.go:102-139`

The repo listing loop pages through all user repos but has no safeguard against GitHub API rate limits. If a user has many repos, the rapid sequential requests could hit rate limits. The code also lacks any delay between pages.

**Fix:** Add rate-limit awareness (check `X-RateLimit-Remaining` headers) or a small delay between paginated requests.

## 7. `godotenv.Load()` error is silently ignored

**File:** `main.go:85`

The error from `godotenv.Load()` is discarded with `_`. If the `.env` file is malformed or missing, the user won't get a helpful message — they'll only see a generic error about `GITHUB_TOKEN` being empty, which may be confusing if the token is present but the file has syntax issues.

**Fix:** Log a warning if `.env` exists but fails to parse, or at least differentiate between "file not found" (which is fine if env vars are set externally) and "file has syntax errors."

---

## TUI Refactor Bugs (found during code review)

### 8. `buildPendingOps()` discarded user-provided rename names

**File:** `tui/select.go` (original version)

`buildPendingOps()` set `m.PendingOps = nil` and then tried to look up rename names from `m.PendingOps` — which it had just cleared. Every rename would fall back to the original repo name, making all renames no-ops.

**Fix:** Added `RenameMap map[string]string` to the model. `buildPendingOps()` now reads new names from `RenameMap` instead of searching `PendingOps`.

### 9. Tab bar highlighting was reversed

**Files:** `tui/select.go`, `tui/confirm.go`

The Select screen highlighted "Confirm" and the Confirm screen highlighted "Select" — the opposite of what the spec requires.

**Fix:** Select screen now uses `ActiveTab.Render("Select") + InactiveTab.Render("Confirm")`. Confirm screen now uses `InactiveTab.Render("Select") + ActiveTab.Render("Confirm")`.

### 10. Viewport initialized in `View()` — mutations lost

**File:** `tui/confirm.go` (original version)

`viewExecuting()` had a value receiver and initialized `m.Viewport` / `m.ViewportReady` on a copy that was never returned. The viewport was recreated on every render and never scrolled because `ViewportReady` was always false in the actual model.

**Fix:** Viewport is now initialized in `updateConfirm()` (which returns the modified model) and `refreshViewport()` is a pointer-receiver helper called from `updateExecuting()`.

### 11. Execution screen had no in-progress indicator

**File:** `tui/confirm.go` (original version)

The spec requires a spinner + "Deleting repo1..." indicator for the currently running operation. The implementation only rendered completed results, so the viewport was blank while the first API call was in flight.

**Fix:** `formatExecResults()` now renders a "... Renaming" or "... Deleting" line for the current in-progress operation above the completed results.