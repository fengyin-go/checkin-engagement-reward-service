package service

import "checkin/internal/model"

// recordCheckinDates 提取签到记录的日期集合，用于连续天数计算。
// 补签记录同样代表该日已签到，应参与连续天数 bridging，故一并纳入。
func recordCheckinDates(records []*model.Checkin) []string {
	dates := make([]string, 0, len(records))
	for _, c := range records {
		dates = append(dates, c.Date)
	}
	return dates
}

func recordUniqueDates(records []*model.Checkin) map[string]bool {
	days := map[string]bool{}
	for _, c := range records {
		days[c.ID] = true
	}
	return days
}
