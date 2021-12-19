package util

import (
	"testing"
)

func TestToUpperCamelCase(t *testing.T) {
	t.Log(ToUpperCamelCase("test_test"))
}

func TestToLowerCamelCase(t *testing.T) {
	t.Log(ToLowerCamelCase(ToUpperCamelCase("test_test")))
}
