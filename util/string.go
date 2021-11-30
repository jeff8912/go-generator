package util

import (
	"strings"
	"unicode"
)

// 转成驼峰
func ToUpperCamelCase(s string) string {
	strs := strings.Split(s, "_")
	result := ""
	for _, str := range strs {
		r := []rune(str)
		result += string(unicode.ToUpper(r[0])) + string(r[1:])
	}
	return result
}

// 转成小驼峰
func ToLowerCamelCase(s string) string {
	result := ""
	r := []rune(s)
	result += string(unicode.ToLower(r[0])) + string(r[1:])
	return result
}
