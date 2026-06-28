package util

import (
	"cmp"
	"slices"
)

func Sort_files(files []string) []string {
	slices.SortFunc(files, func(a, b string) int {
		return cmp.Compare(a, b)
	})
	return files
}
