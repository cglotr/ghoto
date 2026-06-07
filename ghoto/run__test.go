package ghoto

import "testing"

func Test__Run(t *testing.T) {
	Ghoto__new().Run("../testfile/", "Album")
}
