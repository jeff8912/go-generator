package config

import (
	"fmt"
	"testing"
)

func TestGetValue(t *testing.T) {
	fmt.Println(GetValue("system", "server_address"))
}

func TestGetSections(t *testing.T) {
	fmt.Println(GetSections())
}

func TestInt(t *testing.T) {
	fmt.Println(GetDefaultConfig().Int("system", "server_address"))
}
