package diff

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetDiff(t *testing.T) {
	diff := GetDiff(
		[]string{
			"Alice",
			"Bob",
			"Charlie",
			"Daniel",
		},
		[]string{
			"Alice",
			"Bob",
			"Edward",
			"Fredrik",
		},
	)

	expectedDiff := Diff{
		Added: []string{
			"Edward",
			"Fredrik",
		},
		Changed: []string{
			"Alice",
			"Bob",
		},
		Deleted: []string{
			"Charlie",
			"Daniel",
		},
	}

	if !cmp.Equal(diff, expectedDiff) {
		t.Errorf("GetDiff did not return expected value, got=%s", cmp.Diff(expectedDiff, diff))
	}
}
