package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dstotijn/go-notion"
)

type Notion interface {
	QueryDatabase(ctx context.Context, id string, query *notion.DatabaseQuery) (notion.DatabaseQueryResponse, error)
}

type Client struct {
	Notion Notion
}

func (client Client) GetDatabaseRows(databaseId string) (Pages, error) {
	response, err := client.Notion.QueryDatabase(context.Background(), databaseId, nil)

	if err != nil {
		return nil, err
	}

	rows := make(Pages)
	processPages(&response.Results, &rows)

	for response.HasMore {
		response, err = client.Notion.QueryDatabase(context.Background(), databaseId, &notion.DatabaseQuery{
			StartCursor: *response.NextCursor,
		})

		if err != nil {
			return rows, err
		}

		processPages(&response.Results, &rows)
	}

	return rows, nil
}

func processPages(notionPages *[]notion.Page, pages *Pages) {
	for _, notionPage := range *notionPages {
		var row databaseRow

		bytes, err := json.Marshal(notionPage.Properties)

		if err != nil {
			fmt.Println(fmt.Errorf("failed to encode properties for page %s", notionPage.ID))

			continue
		}

		err = json.Unmarshal(bytes, &row)

		if err != nil {
			fmt.Println(fmt.Errorf("failed to decode properies for page %s", notionPage.ID))

			continue
		}

		(*pages)[row.GetSlug()] = page{
			ID: notionPage.ID,
			DB: row,
		}
	}
}
