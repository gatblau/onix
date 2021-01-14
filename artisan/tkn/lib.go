package tkn

import (
	"strings"
)

// encode strings to be used in tekton pipelines names
func encode(value string) string {
	length := 30
	value = strings.ToLower(value)
	value = strings.Replace(value, " ", "-", -1)
	if len(value) > length {
		value = value[0:length]
	}
	return value
}
