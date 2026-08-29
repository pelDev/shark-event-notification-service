package templates

import "strings"

func wrapBase64(s string, lineLen int) string {
	var buf strings.Builder
	for len(s) > 0 {
		chunk := lineLen
		if chunk > len(s) {
			chunk = len(s)
		}
		buf.WriteString(s[:chunk])
		buf.WriteString("\r\n")
		s = s[chunk:]
	}
	return buf.String()
}
