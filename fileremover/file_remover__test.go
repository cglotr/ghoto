package fileremover_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cglotr/ghoto/fileremover"
	"github.com/stretchr/testify/assert"
)

func Test__File_remover__dummy(t *testing.T) {
	r := fileremover.File_remover__dummy__new()
	err := r.Remove("some/file.jpg")
	assert.NoError(t, err)
}

func Test__File_remover__impl(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.jpg")
	f, _ := os.Create(path)
	f.Close()

	r := fileremover.File_remover__impl__new()
	err := r.Remove(path)

	assert.NoError(t, err)
	assert.NoFileExists(t, path)
}
