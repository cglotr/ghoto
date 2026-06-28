package ghoto_test

import (
	"testing"

	"github.com/cglotr/ghoto/ghoto"
	"github.com/stretchr/testify/assert"
)

func Test__Run(t *testing.T) {
	g := ghoto.Ghoto__new()
	err := g.Run("./testfile", "Test")
	assert.Nil(t, err)
}
