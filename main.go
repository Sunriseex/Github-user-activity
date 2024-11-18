package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
}

func fetchGitHubActivity(username string) {
	if username == "" {
		fmt.Println("Please provide a GitHub username.")
		return
	}

	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("User-Agent", "github-user-activity-app")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return

	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 { // 404 is StatusNotFound
		fmt.Printf("Error: Github user '%s' not found. \n", username)
		return
	} else if resp.StatusCode != 200 { // 200 is StatusOK
		fmt.Printf("Error: Failed to fetch data (HTTP %d).\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("err reading response body: %v\n", err)
		return
	}
	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		fmt.Printf("error parsing JSON: %v\n", err)
		return
	}

	fmt.Printf("Recent activity for '%s':\n\n", username)
	for i, event := range events {
		if i >= 10 { // output limit 10 events
			break
		}
		fmt.Printf("- %s in %s\n", event.Type, event.Repo.Name)
	}
}
func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: github-activity <username>")
		return
	}

	username := os.Args[1]
	fetchGitHubActivity(username)
}
