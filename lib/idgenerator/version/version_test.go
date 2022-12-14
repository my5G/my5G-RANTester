package version_test

import (
	"my5G-RANTester/lib/idgenerator/version"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	assert.Equal(t, "2020-05-25-01", version.GetVersion())
}
