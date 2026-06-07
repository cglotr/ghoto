package util

import "strings"

func Filter_photo_files(files []string) []string {
	return filter_files(files, []string{"mp4", "jpg"})
}

func Filter_non_photo_files(files []string) []string {
	return filter_files(files, []string{"lrv", "dng"})
}

func filter_files(files []string, ext []string) []string {
	photo_files := []string{}
	supported_ext := ext

	for _, file := range files {
		words := strings.Split(file, "/")
		if len(words) < 1 {
			continue
		}

		filename := words[len(words)-1]
		words = strings.Split(filename, ".")
		ext := strings.ToLower(words[len(words)-1])
		is_supported := false
		for _, s := range supported_ext {
			if s == ext {
				is_supported = true
			}
		}
		if !is_supported {
			continue
		}

		if filename[0] == '.' {
			continue
		}

		photo_files = append(photo_files, file)
	}

	return photo_files
}
