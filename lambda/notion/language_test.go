package notion

import (
	"fmt"
	"github.com/dstotijn/go-notion"
	"testing"
)

func TestGetColorForLanguage(t *testing.T) {
	langColors := map[string]notion.Color{
		"php":        notion.ColorBlue,
		"typescript": notion.ColorBlue,
		"python":     notion.ColorBlue,
		"go":         notion.ColorBlue,
		"dockerfile": notion.ColorBlue,
		"javascript": notion.ColorYellow,
		"rust":       notion.ColorBrown,
		"c":          notion.ColorBrown,
		"c++":        notion.ColorBrown,
		"shell":      notion.ColorGreen,
		"c#":         notion.ColorGreen,
		"vue":        notion.ColorGreen,
		"html":       notion.ColorGreen,
		"haskell":    notion.ColorRed,
		"swift":      notion.ColorRed,
		"elixir":     notion.ColorPurple,
		"chef":       notion.ColorGray,
		"piet":       notion.ColorGray,
	}

	for lang, color := range langColors {
		language := lang
		expectedColor := color

		t.Run(fmt.Sprintf("%s = %s", lang, color), func(t *testing.T) {
			t.Parallel()

			actualColor := GetColorForLanguage(language)

			if actualColor != expectedColor {
				t.Errorf("Expected GetColorForLanguage(%s) to give %s, got %s", language, expectedColor, actualColor)
			}
		})
	}
}
