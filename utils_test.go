package goose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_StringInSlice_Success(t *testing.T) {

	list := []string{"qwerty", "asdf"}
	result := stringInSlice("qwerty", list)

	assert.True(t, result, "unexpected string in slice result")
}

func Test_StringInSlice_Fail(t *testing.T) {

	list := []string{"qwerty", "asdf"}
	result := stringInSlice("test", list)

	assert.False(t, result, "unexpected string in slice result")
}
