## Carpet Cleaner 🧹

A CLI tool for bulk-deleting GitHub repositories. Lists all repos in your account, lets you pick which ones to delete, and deletes them via the GitHub API.

Future plans include additional operations (rename, archive, etc.) and a TUI interface.

## Prerequisites

- **Go 1.25+** installed
- A **GitHub Personal Access Token** with the `delete_repo` scope enabled

## Setup

1. Clone the repository:
   ```sh
   git clone <repo-url>
   cd carpet-cleaner
   ```

2. Create a `.env` file in the project root:
   ```sh
   cat > .env << 'EOF'
   GITHUB_TOKEN=ghp_your_token_here
   GITHUB_OWNER=your_github_username
   EOF
   ```
   - `GITHUB_TOKEN` — a GitHub Personal Access Token with the `delete_repo` scope
   - `GITHUB_OWNER` — your GitHub username or organization name (used in the API endpoint for deletion)

3. Alternatively, set these as environment variables instead of using `.env`:
   ```sh
   export GITHUB_TOKEN=ghp_your_token_here
   export GITHUB_OWNER=your_github_username
   ```

## Running

```sh
go run main.go
```

Or build and run the binary:

```sh
go build -o carpet-cleaner
./carpet-cleaner
```

## Usage

When you run the tool, it will:

1. **Fetch all your repos** — paginates through the GitHub API and displays each repo with its visibility and description:
   ```
   Fetching repositories...
   Found 12 repositories

   [Private] my-secret-project
       Description: A cool secret thing

   [Public] open-source-lib
       Description: An open source library
   ```

2. **Prompt for deletion** — enter the names of the repos you want to delete, separated by commas:
   ```
   Enter the repos you want to delete separated by commas (e.g. repo1, repo2):
   > my-secret-project, old-repo
   ```

3. **Confirm** — type `YES` (exactly) to confirm deletion:
   ```
   WARNING: You are about to delete 2 repositories.
   Type 'YES' to confirm: YES
   ```

4. **Execute** — the tool deletes each repo and reports success or failure:
   ```
   Attempting to delete repo myuser/my-secret-project... Successfully deleted repository: myuser/my-secret-project
   Attempting to delete repo myuser/old-repo... Successfully deleted repository: myuser/old-repo
   ```

## Creating a GitHub Personal Access Token

1. Go to [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. Click **Generate new token** (classic)
3. Give it a description (e.g., "carpet-cleaner")
4. Under scopes, check **`delete_repo`** — this is the minimum required scope
5. Click **Generate token** and copy the value into your `.env` file

> **Warning:** The `delete_repo` scope grants full delete access to your repositories. Keep your token secure and never commit it to version control.

## TODO

- [ ] Other Operations Aside from Delete
  - [ ] Rename repos (function exists but is not wired up and has bugs)
  - [ ] Update repo descriptions
  - [ ] Archive/unarchive repos
  - [ ] Change repo visibility (public/private)
- [ ] TUI interface for browsing and selecting repos
- [ ] Command-line flags to select operation (delete, rename, etc.) instead of hardcoded flow
- [ ] Validate user-entered repo names against fetched repos before acting
- [ ] Rate-limit awareness when paginating through repos
- [ ] Better error handling for `.env` loading