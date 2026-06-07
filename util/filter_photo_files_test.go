package util_test

import (
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
