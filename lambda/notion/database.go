package notion

func (row *databaseRow) GetSlug() string {
	slug := ""

	for _, titlePart := range row.Name.Title {
		slug += titlePart.PlainText
	}

	return slug
}

func (row *databaseRow) GetPlainTextName() string {
	return row.Name.Title[0].PlainText
}
