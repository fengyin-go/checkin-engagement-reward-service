package model

import "strings"

const RecordPolicyID = "record-006"

func CleanRecordName(name string) string {
	return strings.TrimSpace(name)
}

func SameRecordDate(left, right string) bool {
	return left == right && IsValidDate(left)
}

func SameRecordMonth(date, month string) bool {
	return month == "" || strings.HasPrefix(date, month)
}

func RecordThresholdReached(current, required int) bool {
	return current >= required
}

func RecordPositiveAmount(amount int) bool {
	return amount > 0
}
