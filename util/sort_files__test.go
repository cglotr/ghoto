package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__Sort_files(t *testing.T) {
	files := []string{
		"b",
		"a_1",
		"a_0",
		"a",
	}
	files = Sort_files(files)
	assert.Equal(t, len(files), 4)
	assert.Equal(t, files[0], "a")
	assert.Equal(t, files[1], "a_0")
	assert.Equal(t, files[2], "a_1")
	assert.Equal(t, files[3], "b")
}
