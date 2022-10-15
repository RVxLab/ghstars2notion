package notion

import "github.com/dstotijn/go-notion"

type DatabaseRow struct {
	Name            *NameField            `json:"Name"`
	Description     *DescriptionField     `json:"Description"`
	PrimaryLanguage *PrimaryLanguageField `json:"Primary Language"`
	URL             *URLField             `json:"URL"`
}

type NameField struct {
	Id    string             `json:"id"`
	Title []*notion.RichText `json:"title"`
	Type  string             `json:"type"`
}

type DescriptionField struct {
	Id       string             `json:"id"`
	RichText []*notion.RichText `json:"rich_text"`
	Type     string             `json:"type"`
}

type URLField struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type PrimaryLanguageField struct {
	Id     string                `json:"id"`
	Select *notion.SelectOptions `json:"select"`
	Type   string                `json:"type"`
}
