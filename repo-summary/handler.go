package function

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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

	token := string(tokenData)
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	updatedRepos, _, err := client.Repositories.List(ctx, "", &github.RepositoryListOptions{
		Visibility:  "public",
		Type:        "owner",
		Sort:        "updated",
		Affiliation: "owner",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 10},
	})

	if err != nil {
		log.Printf("Unable to list repositories: %v", err)
		return fmt.Sprintf("Unable to list repos: %v", err)
	}

	summaries := []RepoSummary{}
	for _, repo := range updatedRepos {
		summaries = append(summaries, RepoSummary{
			FullName: repo.GetFullName(),
			Stars:    repo.GetStargazersCount(),
			Issues:   repo.GetOpenIssuesCount(),
			Watchers: repo.GetWatchersCount(),
		})

	}

	out, _ := json.Marshal(summaries)
	return string(out)
}

type RepoSummary struct {
	FullName string
	Stars    int
	Issues   int
	Watchers int
}
