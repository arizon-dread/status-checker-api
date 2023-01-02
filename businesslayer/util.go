package businesslayer

import (
	"strings"
)

func getContentType(s string) string {
	if s != "" {
		if strings.Contains(s, "{") && strings.Contains(s, "\":") {
			return "application/json"
		} else if strings.Contains(s, "</") {
			return "text/xml"
		}
	}
	return "text/plain"
}
