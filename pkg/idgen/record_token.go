package idgen

import "strings"

func RecordToken(parts ...string) string {
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			cleaned = append(cleaned, part)
		}
	}
	if len(cleaned) == 0 {
		return "record-002"
	}
	return strings.Join(cleaned, "-")
}
