package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	goNotion "github.com/dstotijn/go-notion"
	goGithub "github.com/google/go-github/v47/github"
	"notion-sync/lambda/config"
	"notion-sync/lambda/github"
	"notion-sync/lambda/notion"
)

func HandleRequest() (string, error) {
	config := config.GetConfig()

	notionClient := goNotion.NewClient(config.NotionApiKey)

	githubClient := github.Client{
		Starred: goGithub.NewClient(nil).Activity,
	}

	rows, err := notion.GetDatabaseRows(notionClient, config.NotionDatabaseId)

	if rows == nil && err != nil {
		return "Error fetching database", err
	}

	fmt.Printf("Got %d rows\n", len(rows))

	repos, err := githubClient.GetStarredRepos(config.GithubUser)

	if repos == nil && err != nil {
		return "Error fetching starred repos from GitHub", err
	}

	fmt.Printf("Got %d repos\n", len(repos))

	for _, repo := range repos {
		if rows[repo.Name] != nil {
			fmt.Printf("Skipping creation of entry for repo %s as it already exists\n", repo.Name)

			continue
		}

		databasePageProperties := goNotion.DatabasePageProperties{
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

		newPage, err := notionClient.CreatePage(context.Background(), goNotion.CreatePageParams{
			ParentID:               config.NotionDatabaseId,
			ParentType:             goNotion.ParentTypeDatabase,
			DatabasePageProperties: &databasePageProperties,
		})

		if err != nil {
			fmt.Println(fmt.Errorf("failed to create page for repo %s\n", repo.Name))
			fmt.Println(err.Error())

			continue
		}

		fmt.Printf("Created page for repo %s, ID = %s\n", repo.Name, newPage.ID)
	}

	return "Successfully updated Notion database", nil
}

func main() {
	lambda.Start(HandleRequest)
}
