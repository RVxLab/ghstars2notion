package github

import (
	"context"
	"fmt"
	"github.com/google/go-github/v47/github"
)

const reposPerPage = 30

type starredRepo struct {
	Name            string
	PrimaryLanguage string
	Description     string
	URL             string
}

func GetStarredRepos(client *github.Client, username string) ([]*starredRepo, error) {
	repos, response, err := client.Activity.ListStarred(context.Background(), username, &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			PerPage: reposPerPage,
		},
	})

	if err != nil {
		fmt.Printf("Remaining: %d, Limit: %d, Reset: %s", response.Rate.Remaining, response.Rate.Limit, response.Rate.Reset)

		return nil, err
	}

	var starredRepoSlices [][]*starredRepo
	starredRepoSlices = append(starredRepoSlices, processRepos(repos))

	for response.NextPage > 0 {
		repos, response, err = client.Activity.ListStarred(context.Background(), username, &github.ActivityListStarredOptions{
			ListOptions: github.ListOptions{
				PerPage: reposPerPage,
				Page:    response.NextPage,
			},
		})

		if err != nil {
			return nil, err
		}

		starredRepoSlices = append(starredRepoSlices, processRepos(repos))
	}

	var starredRepos []*starredRepo

	for _, starredRepoSlice := range starredRepoSlices {
		starredRepos = append(starredRepos, starredRepoSlice...)
	}

	return starredRepos, nil
}

func processRepos(repos []*github.StarredRepository) []*starredRepo {
	newRepos := make([]*starredRepo, reposPerPage)

	for i, repo := range repos {
		actualRepo := repo.GetRepository()

		language := actualRepo.GetLanguage()

		if language == "" {
			language = "Unknown"
		}

		newRepos[i] = &starredRepo{
			Name:            actualRepo.GetFullName(),
			Description:     actualRepo.GetDescription(),
			PrimaryLanguage: language,
			URL:             actualRepo.GetSVNURL(),
		}
	}

	return newRepos
}
