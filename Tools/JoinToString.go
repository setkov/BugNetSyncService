package Tools

import (
	"strconv"
	"strings"
)

// Convert int slice to string with separator
func JoinToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}
