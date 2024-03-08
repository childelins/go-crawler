package split

import "strings"

func Split(s, sep string) []string {
	var result []string
	i := strings.Index(s, sep)
	for i > -1 {
		result = append(result, s[:i])
		s = s[i+len(sep):]
		i = strings.Index(s, sep)
	}

	// fix with test
	if len(s) > 0 {
		return append(result, s)
	}

	return result
}
