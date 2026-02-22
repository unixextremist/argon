package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	apiURL = "https://api.github.com/search/repositories"
	agent  = "argon-pkg-manager"
	limit  = 10
)

type githubResp struct {
	Items []struct {
		FullName        string `json:"full_name"`
		Description     string `json:"description"`
		StargazersCount int64  `json:"stargazers_count"`
		ForksCount      int64  `json:"forks_count"`
	} `json:"items"`
	Message string `json:"message"`
}

func Search(query string) {
	if query == "" {
		fmt.Fprintln(os.Stderr, "Error: empty search query")
		return
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse URL: %v\n", err)
		return
	}
	
	q := u.Query()
	q.Set("q", query)
	q.Set("sort", "stars")
	q.Set("order", "desc")
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build request: %v\n", err)
		return
	}
	req.Header.Set("User-Agent", agent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to access GitHub API: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Fprintf(os.Stderr, "GitHub API error: %d - %s\n", resp.StatusCode, strings.TrimSpace(string(body)))
		return
	}

	var gr githubResp
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse GitHub response: %v\n", err)
		return
	}

	if len(gr.Items) == 0 {
		if gr.Message != "" {
			fmt.Fprintf(os.Stderr, "GitHub says: %s\n", gr.Message)
		} else {
			fmt.Fprintln(os.Stderr, "No results found")
		}
		return
	}

	for i := 0; i < len(gr.Items) && i < limit; i++ {
		it := gr.Items[i]
		desc := it.Description
		if len(desc) > 50 {
			desc = desc[:47] + "..."
		}
		fmt.Printf("[%s] ⭐ %d 🍴 %d\n", it.FullName, it.StargazersCount, it.ForksCount)
		if desc != "" {
			fmt.Printf("  %s\n", desc)
		}
	}
}
