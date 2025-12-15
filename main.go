package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Repository struct {
	Name        string `json:"name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
}

func RenameRepo(owner, oldName, newName, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, oldName)

	payload := map[string]string {
		"name": newName
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer " + token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to update repo. Status: %d, Response: %s", resp.StatusCode, string(body))
}

func DeleteRepo(owner, repo, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		fmt.Printf("Successfully deleted repository: %s/%s\n", owner, repo)
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("failed to delete repo (Status: %d): %s", resp.StatusCode, string(body))
}

func main() {
	_ = godotenv.Load()

	owner := os.Getenv("GITHUB_OWNER")
	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		log.Fatal("Error: GITHUB_TOKEN is required.")
	}
	if owner == "" {
		log.Fatal("Error: GITHUB_OWNER is required in .env for deletion to work.")
	}

	var allRepos []Repository
	page := 1

	fmt.Println("Fetching repositories...")

	for {
		url := fmt.Sprintf("https://api.github.com/user/repos?per_page=100&page=%d&type=all", page)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatalf("Network error: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close() // Close before exiting
			log.Fatalf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}

		var repos []Repository
		if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
			resp.Body.Close() // Close before crashing
			log.Fatal(err)
		}

		resp.Body.Close()

		if len(repos) == 0 {
			break
		}

		allRepos = append(allRepos, repos...)
		page++
	}

	fmt.Printf("Found %d repositories\n\n", len(allRepos))
	for _, r := range allRepos {
		status := "Public"
		if r.Private {
			status = "Private"
		}

		desc := r.Description
		if desc == "" {
			desc = "(No description)"
		}

		fmt.Printf("[%s] %s\n", status, r.Name)
		fmt.Printf("\tDescription: %s\n\n", desc)
	}

	fmt.Println("Enter the repos you want to delete separated by commas (e.g. repo1, repo2):")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println("No repositories entered. Exiting.")
		os.Exit(0)
	}

	repoList := strings.Split(input, ",")
	for i := range repoList {
		repoList[i] = strings.TrimSpace(repoList[i])
	}

	fmt.Printf("WARNING: You are about to delete %d repositories.\n", len(repoList))
	fmt.Print("Type 'YES' to confirm: ")

	confirmStr, _ := reader.ReadString('\n')
	confirmStr = strings.TrimSpace(confirmStr)

	if confirmStr != "YES" {
		fmt.Println("Confirmation failed. Exiting.")
		return
	}

	for _, repo := range repoList {
		if repo == "" {
			continue
		}
		fmt.Printf("Attempting to delete repo %s/%s... ", owner, repo)
		err := DeleteRepo(owner, repo, token)
		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
		}
	}
}
