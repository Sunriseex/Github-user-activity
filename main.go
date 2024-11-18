package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	cacheFileName = "github_activity_cache.json"
	cacheExpiry   = time.Hour
)

type Event struct {
	Type string `json:"type"`
	Repo struct {
		Name string `json:"name"`
	} `json:"repo"`
	CreatedAt string `json:"created_at"`
}

type CacheEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Events    []Event   `json:"events"`
}

var cache = make(map[string]CacheEntry)

func loadCache() {
	data, err := os.ReadFile(cacheFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Printf("Error reading cache file: %v\n", err)
		return
	}
	if err := json.Unmarshal(data, &cache); err != nil {
		fmt.Printf("Error parsing cache file: %v\n", err)
	}
}

func saveCache() {
	data, err := json.MarshalIndent(cache, "", " ")
	if err != nil {
		fmt.Printf("Error serializing cache: %v\n", err)
		return
	}
	if err := os.WriteFile(cacheFileName, data, 0644); err != nil {
		fmt.Printf("Error writing cache file: %v\n", err)
	}
}

func getCachedActivity(username string) ([]Event, bool) {
	entry, exists := cache[username]
	if !exists {
		return nil, false
	}
	if time.Since(entry.Timestamp) > cacheExpiry {
		return nil, false
	}
	return entry.Events, true
}

func fetchGitHubActivity(username string, filterType string) {
	if events, found := getCachedActivity(username); found {
		fmt.Printf("Using cached data for '%s':\n\n", username)
		displayActivity(events, filterType)
		return
	}
	fmt.Printf("Fetching data for '%s' from GitHub API...\n", username)

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

	cache[username] = CacheEntry{
		Timestamp: time.Now(),
		Events:    events,
	}
	saveCache()
	displayActivity(events, filterType)
}

func displayActivity(events []Event, filterType string) {
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
	loadCache()

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
