package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Stats struct {
	Commits      int `json:"commits"`
	PullRequests int `json:"pullRequests"`
	Contributed  int `json:"contributed"`
	Stars        int `json:"stars"`
}

type GraphQLResponse struct {
	Data struct {
		User struct {
			Commits struct {
				Total int `json:"totalCommitContributions"`
			} `json:"contributionsCollection"`
			PullRequests struct {
				Total int `json:"totalCount"`
			} `json:"pullRequests"`
			Contributed struct {
				Total int `json:"totalCount"`
			} `json:"repositoriesContributedTo"`
			Stars struct {
				Nodes []struct {
					Stars struct {
						Total int `json:"totalCount"`
					} `json:"stargazers"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
}

func GetGithubStats() (*Stats, error) {
	token := os.Getenv("GH_ACCESS_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GH_ACCESS_TOKEN environment variable not set")
	}

	yearAgo := time.Now().AddDate(-1, 0, 0).Format(time.RFC3339)

	query := fmt.Sprintf(`
	{
		user(login: "xqsit94") {
			contributionsCollection(from: "%s", to: "%s") {
				totalCommitContributions
			}
			pullRequests {
				totalCount
			}
			repositoriesContributedTo(contributionTypes: [COMMIT, ISSUE, PULL_REQUEST, REPOSITORY]) {
				totalCount
			}
			repositories(ownerAffiliations: OWNER, first: 100, orderBy: {field: STARGAZERS, direction: DESC}) {
				nodes {
					stargazers {
						totalCount
					}
				}
			}
		}
	}`, yearAgo, time.Now().Format(time.RFC3339))

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", strings.NewReader(fmt.Sprintf(`{"query": %q}`, query)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var graphQLResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphQLResp); err != nil {
		return nil, err
	}

	totalStars := 0
	for _, node := range graphQLResp.Data.User.Stars.Nodes {
		totalStars += node.Stars.Total
	}

	return &Stats{
		Commits:      graphQLResp.Data.User.Commits.Total,
		PullRequests: graphQLResp.Data.User.PullRequests.Total,
		Contributed:  graphQLResp.Data.User.Contributed.Total,
		Stars:        totalStars,
	}, nil
}
