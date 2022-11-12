package notion

func (row *databaseRow) GetSlug() string {
	slug := ""

	for _, titlePart := range row.Name.Title {
		slug += titlePart.PlainText
	}

	return slug
}
