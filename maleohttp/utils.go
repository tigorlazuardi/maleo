package maleohttp

import (
	"encoding/json"
	"strings"
)

func isHumanReadable(contentType string) bool {
	return contentType == "" ||
		strings.Contains(contentType, "text/") ||
		strings.Contains(contentType, "application/json") ||
		strings.Contains(contentType, "application/xml")
}

func isJson(b []byte) bool {
	var js json.RawMessage
	return json.Unmarshal(b, &js) == nil
}
