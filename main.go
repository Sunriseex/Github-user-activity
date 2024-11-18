package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt string `json:"created_at"`
}

func fetchGitHubActivity(username string, filterType string) {
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

	if len(events) == 0 {
		fmt.Printf("No recent activities found for user '%s'\n", username)
		return
	}

	fmt.Printf("Recent activity for '%s':\n\n", username)
	count := 0
	for _, event := range events {
		if filterType != "" && event.Type != filterType {
			continue
		}
		count++
		if count > 10 {
			break
		}
		createdAt, err := time.Parse(time.RFC3339, event.CreatedAt)
		if err != nil {
			fmt.Printf("Error parsing event %v\n", err)
			continue
		}
		fmt.Printf("- %s in %s on %s\n", event.Type, event.Repo.Name, createdAt.Format("2006-01-02 15:04:05"))
	}

	if count == 0 {
		fmt.Printf("No events found matching the filter '%s'.\n", filterType)
	}
}
func main() {

	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Println("Usage: github-activity <username> [event-type]")
		return
	}

	username := os.Args[1]
	var filterType string
	if len(os.Args) == 3 {
		filterType = os.Args[2]
	}
	fetchGitHubActivity(username, filterType)
}
