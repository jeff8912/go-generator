package util

import (
	"strings"
	"unicode"
)

// ToUpperCamelCase 转成驼峰
func ToUpperCamelCase(s string) string {
	splits := strings.Split(s, "_")
	result := ""
	for _, str := range splits {
		r := []rune(str)
		result += string(unicode.ToUpper(r[0])) + string(r[1:])
	}
	return result
}

// ToLowerCamelCase 大驼峰转成小驼峰
func ToLowerCamelCase(s string) string {
	result := ""
	r := []rune(s)
	result += string(unicode.ToLower(r[0])) + string(r[1:])
	return result
}
