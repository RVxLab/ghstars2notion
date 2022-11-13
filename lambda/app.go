package main

import (
	"context"
	"fmt"
	"g2stars2notion/lambda/config"
	"g2stars2notion/lambda/diff"
	"g2stars2notion/lambda/github"
	"g2stars2notion/lambda/notion"
	"g2stars2notion/lambda/utils"
	goNotion "github.com/dstotijn/go-notion"
	goGithub "github.com/google/go-github/v47/github"
)

func makePagePropertiesPayload(repo github.StarredRepository) goNotion.DatabasePageProperties {
	return goNotion.DatabasePageProperties{
		"Name": goNotion.DatabasePageProperty{
			ID:   "title",
			Type: goNotion.DBPropTypeTitle,
			Title: []goNotion.RichText{
				{
					Type: goNotion.RichTextTypeText,
					Text: &goNotion.Text{
						Content: repo.Name,
					},
				},
			},
		},
		"Primary Language": goNotion.DatabasePageProperty{
			Type: goNotion.DBPropTypeSelect,
			Select: &goNotion.SelectOptions{
				Name:  repo.PrimaryLanguage,
				Color: notion.GetColorForLanguage(repo.PrimaryLanguage),
			},
		},
		"Description": goNotion.DatabasePageProperty{
			Type: goNotion.DBPropTypeRichText,
			RichText: []goNotion.RichText{
				{
					Type: goNotion.RichTextTypeText,
					Text: &goNotion.Text{
						Content: repo.Description,
					},
				},
			},
		},
		"URL": goNotion.DatabasePageProperty{
			Type: goNotion.DBPropTypeURL,
			URL:  &repo.URL,
		},
	}
}

func createEntry(repo github.StarredRepository, notionClient *goNotion.Client, databaseId string) (goNotion.Page, error) {
	databasePageProperties := makePagePropertiesPayload(repo)

	return notionClient.CreatePage(context.Background(), goNotion.CreatePageParams{
		ParentID:               databaseId,
		ParentType:             goNotion.ParentTypeDatabase,
		DatabasePageProperties: &databasePageProperties,
	})
}

func deleteEntry(notionClient *goNotion.Client, pageId string) error {
	_, err := notionClient.UpdatePage(context.Background(), pageId, goNotion.UpdatePageParams{
		Archived: utils.Bool(true),
	})

	return err
}

func run() (string, error) {
	appConfig := config.GetConfig()

	goNotionClient := goNotion.NewClient(appConfig.NotionApiKey)

	notionClient := notion.Client{
		Notion: goNotionClient,
	}

	githubClient := github.Client{
		Starred: goGithub.NewClient(nil).Activity,
	}

	rows, err := notionClient.GetDatabaseRows(appConfig.NotionDatabaseId)

	if rows == nil && err != nil {
		return "Error fetching database", err
	}

	fmt.Printf("Got %d rows\n", len(rows))

	repos, err := githubClient.GetStarredRepos(appConfig.GithubUser)

	if repos == nil && err != nil {
		return "Error fetching starred repos from GitHub", err
	}

	fmt.Printf("Got %d repos\n", len(repos))

	repoDiff := diff.GetDiff(utils.MapGetKeys(rows), utils.MapGetKeys(repos))

	fmt.Printf("Added = %d, Deleted = %d, Unchanged = %d\n", len(repoDiff.Added), len(repoDiff.Deleted), len(repoDiff.Changed))

	for _, repoKey := range repoDiff.Added {
		if repo, ok := repos[repoKey]; ok {
			newPage, err := createEntry(repo, goNotionClient, appConfig.NotionDatabaseId)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Created entry for repo %s, ID = %s\n", repo.Name, newPage.ID)
			}
		} else {
			fmt.Println(fmt.Errorf("could not find repo with key %s", repoKey))
		}
	}

	for _, rowKey := range repoDiff.Deleted {
		if row, ok := rows[rowKey]; ok {
			err := deleteEntry(goNotionClient, row.ID)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Deleted entry %s, ID = %s\n", row.DB.GetPlainTextName(), row.ID)
			}
		} else {
			fmt.Println(fmt.Errorf("could not find row with key %s", rowKey))
		}
	}

	return "Successfully updated Notion database", nil
}
