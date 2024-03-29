package function

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

// Handle a serverless request
func Handle(req []byte) string {

	tokenData, err := os.ReadFile("/var/openfaas/secrets/repo-reader-token")
	if err != nil {
		log.Printf("Unable to read token file: %v", err)
		return "Unable to read token file"
	}

	token := strings.TrimSpace(string(tokenData))
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	updatedRepos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{
		Visibility:  "public",
		Affiliation: "owner",
		// Type:        "owner",
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 10},
	})

	if err != nil {
		log.Printf("Unable to list repositories: %v", err)
		return "Unable to list repositories"
	}

	summaries := []RepoSummary{}
	for _, repo := range updatedRepos {
		summaries = append(summaries, RepoSummary{
			FullName: repo.GetFullName(),
			Stars:    repo.GetStargazersCount(),
			Issues:   repo.GetOpenIssuesCount(),
			Watchers: repo.GetWatchersCount(),
			Updated:  repo.GetUpdatedAt().Time,
			// URL:      repo.GetHTMLURL(),
		})

	}

	out, _ := json.Marshal(summaries)
	return string(out)
}

type RepoSummary struct {
	FullName string    `json:"full_name"`
	URL      string    `json:"url"`
	Stars    int       `json:"stars"`
	Issues   int       `json:"issues"`
	Watchers int       `json:"watchers"`
	Updated  time.Time `json:"updated"`
}
