package notion

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dstotijn/go-notion"
)

func GetDatabaseRows(client *notion.Client, databaseId string) (map[string]*DatabaseRow, error) {
	response, err := client.QueryDatabase(context.Background(), databaseId, nil)

	if err != nil {
		return nil, err
	}

	rows := make(map[string]*DatabaseRow)
	processPages(&response.Results, &rows)

	for response.HasMore {
		response, err = client.QueryDatabase(context.Background(), databaseId, &notion.DatabaseQuery{
			StartCursor: *response.NextCursor,
		})

		if err != nil {
			return rows, err
		}

		processPages(&response.Results, &rows)
	}

	return rows, nil
}

func processPages(notionPages *[]notion.Page, pages *map[string]*DatabaseRow) {
	for _, page := range *notionPages {
		var row DatabaseRow

		bytes, err := json.Marshal(page.Properties)

		if err != nil {
			fmt.Errorf("failed to encode properties for page %s", page.ID)

			continue
		}

		err = json.Unmarshal(bytes, &row)

		if err != nil {
			fmt.Errorf("failed to decode properies for page %s", page.ID)

			continue
		}

		(*pages)[row.GetSlug()] = &row
	}
}
