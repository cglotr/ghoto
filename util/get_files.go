package util

import (
	"io/fs"
	"path/filepath"
	"sort"
)

func Get_files(dir string) []string {
	var err error
	var files []string = []string{}
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic("Get_files__WalkDir: " + err.Error())
	}
	sort.Strings(files)
	return files
}
