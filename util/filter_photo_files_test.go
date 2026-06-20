package util_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cglotr/ghoto/util"
	"github.com/stretchr/testify/assert"
)

func Test__Filter_photo_files(t *testing.T) {
	photo_files := util.Filter_photo_files(util.Get_files("./testfile"))
	assert.Equal(t, 2, len(photo_files))

	assert.Equal(t, "testfile/photo_20260613.jpg", photo_files[0])
	assert.Equal(t, "testfile/video_20260613.mp4", photo_files[1])
}

func Test__Filter_non_photo_files(t *testing.T) {
	non_photo_files := util.Filter_non_photo_files(util.Get_files("./testfile"))
	assert.Equal(t, 2, len(non_photo_files))

	assert.Equal(t, "testfile/photo.dng", non_photo_files[0])
	assert.Equal(t, "testfile/video.lrv", non_photo_files[1])
}

func Test__Filter_photo_files__empty(t *testing.T) {
	result := util.Filter_photo_files([]string{})
	assert.Equal(t, 0, len(result))
}

func Test__Get_files__sorted(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"c.jpg", "a.jpg", "b.jpg"} {
		f, _ := os.Create(filepath.Join(dir, name))
		f.Close()
	}
	files := util.Get_files(dir)
	assert.Equal(t, 3, len(files))
	assert.Equal(t, filepath.Join(dir, "a.jpg"), files[0])
	assert.Equal(t, filepath.Join(dir, "b.jpg"), files[1])
	assert.Equal(t, filepath.Join(dir, "c.jpg"), files[2])
}

func Test__Get_files__nonexistent_dir_panics(t *testing.T) {
	assert.Panics(t, func() {
		util.Get_files("/nonexistent/path/that/does/not/exist")
	})
}
