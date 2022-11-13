package notion

import (
	"context"
	"fmt"
	"g2stars2notion/lambda/utils"
	"github.com/dstotijn/go-notion"
	"github.com/google/go-cmp/cmp"
	"testing"
)

var testPages = []notion.Page{
	{
		ID: "P1",
		Properties: notion.DatabasePageProperties{
			"URL": {
				ID:   "Url1",
				Type: notion.DBPropTypeURL,
				URL:  utils.String("https://example.com/page1"),
			},
			"Description": {
				ID:   "Desc1",
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{
						Type: notion.RichTextTypeText,
						Text: &notion.Text{
							Content: "The first page",
						},
						Annotations: &notion.Annotations{
							Bold:          false,
							Italic:        false,
							Strikethrough: false,
							Underline:     false,
							Code:          false,
							Color:         notion.ColorDefault,
						},
						PlainText: "The first page",
					},
				},
			},
			"Primary Language": {
				ID:   "Lang1",
				Type: notion.DBPropTypeSelect,
				Select: &notion.SelectOptions{
					ID:    "Select1",
					Name:  "JS",
					Color: notion.ColorGreen,
				},
			},
			"Name": {
				ID:   "title",
				Type: notion.DBPropTypeTitle,
				Title: []notion.RichText{
					{
						Type: notion.RichTextTypeText,
						Text: &notion.Text{
							Content: "Page 1",
						},
						Annotations: &notion.Annotations{
							Bold:          false,
							Italic:        false,
							Strikethrough: false,
							Underline:     false,
							Code:          false,
							Color:         notion.ColorDefault,
						},
						PlainText: "Page 1",
					},
				},
			},
		},
	},
	{
		ID: "P2",
		Properties: notion.DatabasePageProperties{
			"URL": {
				ID:   "Url2",
				Type: notion.DBPropTypeURL,
				URL:  utils.String("https://example.com/page2"),
			},
			"Description": {
				ID:   "Desc2",
				Type: notion.DBPropTypeRichText,
				RichText: []notion.RichText{
					{
						Type: notion.RichTextTypeText,
						Text: &notion.Text{
							Content: "The second page",
						},
						Annotations: &notion.Annotations{
							Bold:          false,
							Italic:        false,
							Strikethrough: false,
							Underline:     false,
							Code:          false,
							Color:         notion.ColorDefault,
						},
						PlainText: "The second page",
					},
				},
			},
			"Primary Language": {
				ID:   "Lang2",
				Type: notion.DBPropTypeSelect,
				Select: &notion.SelectOptions{
					ID:    "Select2",
					Name:  "PHP",
					Color: notion.ColorPurple,
				},
			},
			"Name": {
				ID:   "title",
				Type: notion.DBPropTypeTitle,
				Title: []notion.RichText{
					{
						Type: notion.RichTextTypeText,
						Text: &notion.Text{
							Content: "Page 2",
						},
						Annotations: &notion.Annotations{
							Bold:          false,
							Italic:        false,
							Strikethrough: false,
							Underline:     false,
							Code:          false,
							Color:         notion.ColorDefault,
						},
						PlainText: "Page 2",
					},
				},
			},
		},
	},
}

type TestClient struct{}

func (client *TestClient) QueryDatabase(ctx context.Context, id string, query *notion.DatabaseQuery) (notion.DatabaseQueryResponse, error) {
	if query == nil {
		response := notion.DatabaseQueryResponse{
			Results: []notion.Page{
				testPages[0],
			},
			HasMore:    true,
			NextCursor: utils.String("Bla"),
		}

		return response, nil
	}

	response := notion.DatabaseQueryResponse{
		Results: []notion.Page{
			testPages[1],
		},
		HasMore:    false,
		NextCursor: nil,
	}

	return response, nil
}

type InstantlyFailingClient struct{}

func (client *InstantlyFailingClient) QueryDatabase(ctx context.Context, id string, query *notion.DatabaseQuery) (notion.DatabaseQueryResponse, error) {
	response := notion.DatabaseQueryResponse{
		Results:    []notion.Page{},
		HasMore:    false,
		NextCursor: nil,
	}

	return response, fmt.Errorf("the request failed")
}

type AfterFirstCallFailingClient struct{}

func (client *AfterFirstCallFailingClient) QueryDatabase(ctx context.Context, id string, query *notion.DatabaseQuery) (notion.DatabaseQueryResponse, error) {
	if query == nil {
		response := notion.DatabaseQueryResponse{
			Results: []notion.Page{
				testPages[0],
			},
			HasMore:    true,
			NextCursor: utils.String("Cursor to the next page"),
		}

		return response, nil
	}

	response := notion.DatabaseQueryResponse{
		Results:    []notion.Page{},
		HasMore:    false,
		NextCursor: nil,
	}

	return response, fmt.Errorf("hit rate limit")
}

func TestGetDatabaseRows(t *testing.T) {
	var testClient TestClient

	client := Client{
		Notion: &testClient,
	}

	rows, err := client.GetDatabaseRows("test")

	if err != nil {
		t.Errorf("Expected GetDatabaseRows to not error, got %s", err)

		return
	}

	page1Properties := testPages[0].Properties.(notion.DatabasePageProperties)
	page2Properties := testPages[1].Properties.(notion.DatabasePageProperties)

	expectedRows := Pages{
		"Page 1": {
			ID: "P1",
			DB: databaseRow{
				Name: nameField{
					ID:    "title",
					Type:  "title",
					Title: page1Properties["Name"].Title,
				},
				Description: descriptionField{
					ID:       "Desc1",
					Type:     "rich_text",
					RichText: page1Properties["Description"].RichText,
				},
				URL: urlField{
					ID:   "Url1",
					Type: "url",
					Url:  "https://example.com/page1",
				},
				PrimaryLanguage: primaryLanguageField{
					ID:     "Lang1",
					Type:   "select",
					Select: page1Properties["Primary Language"].Select,
				},
			},
		},
		"Page 2": {
			ID: "P2",
			DB: databaseRow{
				Name: nameField{
					ID:    "title",
					Type:  "title",
					Title: page2Properties["Name"].Title,
				},
				Description: descriptionField{
					ID:       "Desc2",
					Type:     "rich_text",
					RichText: page2Properties["Description"].RichText,
				},
				URL: urlField{
					ID:   "Url2",
					Type: "url",
					Url:  "https://example.com/page2",
				},
				PrimaryLanguage: primaryLanguageField{
					ID:     "Lang2",
					Type:   "select",
					Select: page2Properties["Primary Language"].Select,
				},
			},
		},
	}

	if !cmp.Equal(rows, expectedRows) {
		t.Errorf("Unexpected output from GetDatabaseRows, got=%s", cmp.Diff(expectedRows, rows))
	}
}

func TestGetDatabaseRowsFailsInitially(t *testing.T) {
	var failingClient InstantlyFailingClient

	client := Client{
		Notion: &failingClient,
	}

	_, err := client.GetDatabaseRows("test")

	if err == nil {
		t.Error("Expected GetDatabaseRows to error, got nil")
	}
}

func TestGetDatabaseRowsFailsAtSecondCall(t *testing.T) {
	var failingClient AfterFirstCallFailingClient

	client := Client{
		Notion: &failingClient,
	}

	_, err := client.GetDatabaseRows("test")

	if err == nil {
		t.Error("Expected GetDatabaseRows to error, got nil")
	}
}
