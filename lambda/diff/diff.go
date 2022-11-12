package diff

import "notion-sync/lambda/utils"

type Diff struct {
	Added   []string
	Changed []string
	Deleted []string
}

func GetDiff(existing, new []string) Diff {
	diff := Diff{
		Added:   []string{},
		Changed: []string{},
		Deleted: []string{},
	}

	for _, v := range new {
		if utils.Contains[string](existing, v) {
			diff.Changed = append(diff.Changed, v)
		} else {
			diff.Added = append(diff.Added, v)
		}
	}

	for _, v := range existing {
		if !utils.Contains[string](new, v) {
			diff.Deleted = append(diff.Deleted, v)
		}
	}

	return diff
}
