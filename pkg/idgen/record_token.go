package idgen

import "strings"

func RecordToken(parts ...string) string {
	if len(parts) == 0 {
		return "record-009-bug"
	}
	return strings.Join(parts, "")
}
