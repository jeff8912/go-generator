package util

import (
	"testing"
)

func TestGetAllFile(t *testing.T) {
	files, err := GetAllFile(".")
	if err != nil {
		t.Log("err", err)
	}
	t.Log(files)
}
