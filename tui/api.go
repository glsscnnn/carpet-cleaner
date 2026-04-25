package tui

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ghRepo struct {
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func fetchAllRepos(owner, token string) ([]RepoItem, error) {
	var all []RepoItem
	page := 1

	for {
		url := fmt.Sprintf("https://api.github.com/user/repos?per_page=100&page=%d&type=all", page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		setAuthHeaders(req, token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("network error: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("API request failed (status %d): %s", resp.StatusCode, string(body))
		}

		var repos []ghRepo
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		if len(repos) == 0 {
			break
		}

		for _, r := range repos {
			all = append(all, RepoItem{
				Name:    r.Name,
				Owner:   r.Owner.Login,
				Desc:    r.Description,
				Private: r.Private,
			})
		}
		page++
	}

	return all, nil
}

func deleteRepo(owner, repo, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	setAuthHeaders(req, token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to delete repo (status %d): %s", resp.StatusCode, string(body))
}

func renameRepo(owner, oldName, newName, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, oldName)

	payload := map[string]string{"name": newName}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	setAuthHeaders(req, token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to rename repo (status %d): %s", resp.StatusCode, string(body))
}

func setAuthHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}