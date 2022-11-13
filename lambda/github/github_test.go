package github

import (
	"context"
	"fmt"
	"g2stars2notion/lambda/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-github/v47/github"
	"testing"
)

var testRepos = []*github.StarredRepository{
	{
		Repository: &github.Repository{
			FullName:    utils.String("User/Repo1"),
			Description: utils.String("PHP library that does things"),
			Language:    utils.String("PHP"),
			SVNURL:      utils.String("https://example.com/repo1.git"),
		},
	},
	{
		Repository: &github.Repository{
			FullName:    utils.String("User/Repo2"),
			Description: utils.String("Rust CLI app that handles things"),
			Language:    utils.String("Rust"),
			SVNURL:      utils.String("https://example.com/repo2.git"),
		},
	},
	{
		Repository: &github.Repository{
			FullName:    utils.String("User/Repo3"),
			Description: utils.String("JavaScript library that calculates the 20th digit of pi"),
			Language:    utils.String("JavaScript"),
			SVNURL:      utils.String("https://example.com/repo3.git"),
		},
	},
	{
		Repository: &github.Repository{
			FullName:    utils.String("User/Repo4"),
			Description: utils.String(""),
			Language:    utils.String(""),
			SVNURL:      utils.String("https://example.com/repo4.git"),
		},
	},
}

type TestClient struct{}

func (client *TestClient) ListStarred(ctx context.Context, user string, opts *github.ActivityListStarredOptions) (starredRepos []*github.StarredRepository, response *github.Response, error error) {
	var repos []*github.StarredRepository
	var resp github.Response

	if opts.Page == 0 {
		// Initial response
		resp.NextPage = 2
		repos = []*github.StarredRepository{
			testRepos[0],
			testRepos[1],
		}
	} else {
		// Next response
		resp.NextPage = 0
		repos = []*github.StarredRepository{
			testRepos[2],
			testRepos[3],
		}
	}

	return repos, &resp, nil
}

type InstantlyFailingClient struct{}

func (client *InstantlyFailingClient) ListStarred(ctx context.Context, user string, opts *github.ActivityListStarredOptions) (starredRepos []*github.StarredRepository, response *github.Response, error error) {
	return nil, nil, fmt.Errorf("the request failed")
}

type AfterFirstCallFailingClient struct{}

func (client *AfterFirstCallFailingClient) ListStarred(ctx context.Context, user string, opts *github.ActivityListStarredOptions) (starredRepos []*github.StarredRepository, response *github.Response, error error) {
	if opts.Page == 0 {
		repos := []*github.StarredRepository{
			testRepos[0],
		}

		resp := github.Response{
			NextPage: 2,
		}

		return repos, &resp, nil
	}

	return nil, nil, fmt.Errorf("hit rate limit")
}

func TestGetStarredRepos(t *testing.T) {
	var testClient TestClient

	client := Client{
		Starred: &testClient,
	}

	repos, _ := client.GetStarredRepos("Test user")

	expectedRepos := StarredRepositories{
		"User/Repo1": {
			Name:            "User/Repo1",
			Description:     "PHP library that does things",
			PrimaryLanguage: "PHP",
			URL:             "https://example.com/repo1.git",
		},
		"User/Repo2": {
			Name:            "User/Repo2",
			Description:     "Rust CLI app that handles things",
			PrimaryLanguage: "Rust",
			URL:             "https://example.com/repo2.git",
		},
		"User/Repo3": {
			Name:            "User/Repo3",
			Description:     "JavaScript library that calculates the 20th digit of pi",
			PrimaryLanguage: "JavaScript",
			URL:             "https://example.com/repo3.git",
		},
		"User/Repo4": {
			Name:            "User/Repo4",
			Description:     "",
			PrimaryLanguage: "Unknown",
			URL:             "https://example.com/repo4.git",
		},
	}

	if !cmp.Equal(repos, expectedRepos) {
		t.Errorf("Unexpected output from GetStarredRepos, got=%s", cmp.Diff(expectedRepos, repos))
	}
}

func TestGetStarredReposFailsInitially(t *testing.T) {
	var failingClient InstantlyFailingClient

	client := Client{
		Starred: &failingClient,
	}

	_, err := client.GetStarredRepos("Test user")

	if err == nil {
		t.Error("Expected GetStarredRepo to error, got nil")
	}
}

func TestGetStarredReposFailsAtSecondCall(t *testing.T) {
	var failingClient AfterFirstCallFailingClient

	client := Client{
		Starred: &failingClient,
	}

	_, err := client.GetStarredRepos("Test user")

	if err == nil {
		t.Error("Expected GetStarredRepo to error, got nil")
	}
}
