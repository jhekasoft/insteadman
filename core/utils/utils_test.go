package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBinAbsDir(t *testing.T) {
	pathList := map[string]string{
		"/usr/bin/insteadman-gtk": "/usr/bin",
	}

	for executablePath, mustBeDir := range pathList {
		dir, e := BinAbsDir(executablePath)
		assert.NoError(t, e)
		assert.Equal(t, dir, mustBeDir)
	}
}
