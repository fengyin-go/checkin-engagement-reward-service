package service

import "checkin/internal/model"

func recordCheckinDates(records []*model.Checkin) []string {
	dates := make([]string, 0, len(records))
	for _, c := range records {
		if c.Source != model.SourceMakeup {
			dates = append(dates, c.Date)
		}
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
