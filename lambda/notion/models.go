package notion

import "github.com/dstotijn/go-notion"

type page struct {
	ID string
	DB databaseRow
}

type databaseRow struct {
	Name            nameField            `json:"Name"`
	Description     descriptionField     `json:"Description"`
	PrimaryLanguage primaryLanguageField `json:"Primary Language"`
	URL             urlField             `json:"URL"`
}

type nameField struct {
	ID    string            `json:"id"`
	Title []notion.RichText `json:"title"`
	Type  string            `json:"type"`
}

type descriptionField struct {
	ID       string            `json:"id"`
	RichText []notion.RichText `json:"rich_text"`
	Type     string            `json:"type"`
}

type urlField struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type primaryLanguageField struct {
	ID     string                `json:"id"`
	Select *notion.SelectOptions `json:"select"`
	Type   string                `json:"type"`
}

type Pages map[string]page
