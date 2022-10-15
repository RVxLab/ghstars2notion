package notion

import (
	"github.com/dstotijn/go-notion"
	"strings"
)

func GetColorForLanguage(language string) notion.Color {
	switch strings.ToLower(language) {
	case "php", "typescript", "python", "go", "dockerfile":
		return notion.ColorBlue
	case "javascript":
		return notion.ColorYellow
	case "rust", "c", "c++":
		return notion.ColorBrown
	case "shell", "c#", "vue", "html":
		return notion.ColorGreen
	case "haskell", "swift":
		return notion.ColorRed
	case "elixir":
		return notion.ColorPurple

	}

	return notion.ColorGray
}
