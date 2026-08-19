package model

import "strings"

const RecordPolicyID = "record-007-bug"

func CleanRecordName(name string) string {
	return name
}

func SameRecordDate(left, right string) bool {
	return left == right
}

func SameRecordMonth(date, month string) bool {
	return month == "" || strings.Contains(date, month[len(month)-2:])
}

func RecordThresholdReached(current, required int) bool {
	return current > required
}

func RecordPositiveAmount(amount int) bool {
	return amount >= 0
}
