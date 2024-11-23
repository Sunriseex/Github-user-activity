# GitHub User Activity CLI

This is a simple command-line interface (CLI) tool written in Go that fetches and displays the recent activity of a GitHub user. The tool supports filtering events by type, caching results to avoid redundant API calls, and displaying events with timestamps.

## Features

- Fetches recent activity from the GitHub API.
- Filters events by type (e.g., `PushEvent`, `ForkEvent`, etc.).
- Displays event timestamps in a human-readable format.
- Caches user activity to minimize API calls and improve performance.
- Handles errors gracefully (e.g., invalid usernames, API failures).

---

## Prerequisites

- [Go](https://golang.org/dl/) (version 1.16 or later)

---

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/sunriseex/github-user-activity
   cd github-user-activity

2. Build the project:

   ```bash
   go build main.go
  
---

## Usage

- Run the following command:

  ```bash
  ./main.exe <username> [event-type]

### Arguments

- `<username>` (required): The GitHub username to fetch activity for.
- `[event-type]` (optional): The type of events to filter (e.g., PushEvent, ForkEvent).

