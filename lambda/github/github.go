package github

import (
	"context"
	"github.com/google/go-github/v47/github"
)

const reposPerPage = 30

type starredRepository struct {
	Name            string
	PrimaryLanguage string
	Description     string
	URL             string
}

type starredRepositories map[string]starredRepository

type ListsStarredRepos interface {
	ListStarred(ctx context.Context, user string, opts *github.ActivityListStarredOptions) ([]*github.StarredRepository, *github.Response, error)
}

type Client struct {
	Starred ListsStarredRepos
}

func (client *Client) GetStarredRepos(username string) (starredRepositories, error) {
	repos, response, err := client.Starred.ListStarred(context.Background(), username, &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			PerPage: reposPerPage,
		},
	})

	if err != nil {
		return nil, err
	}

	starredRepos := make(starredRepositories)

	processRepos(repos, &starredRepos)

	for response.NextPage > 0 {
		repos, response, err = client.Starred.ListStarred(context.Background(), username, &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				PerPage: reposPerPage,
				Page:    response.NextPage,
			},
		})

		if err != nil {
			return nil, err
		}

		processRepos(repos, &starredRepos)
	}

	return starredRepos, nil
}

func processRepos(repos []*github.StarredRepository, starredRepos *starredRepositories) {
	for _, repo := range repos {
		actualRepo := repo.GetRepository()

		language := actualRepo.GetLanguage()

		if language == "" {
			language = "Unknown"
		}

		repoName := actualRepo.GetFullName()

		(*starredRepos)[repoName] = starredRepository{
			Name:            repoName,
			Description:     actualRepo.GetDescription(),
			PrimaryLanguage: language,
			URL:             actualRepo.GetSVNURL(),
		}
	}
}
