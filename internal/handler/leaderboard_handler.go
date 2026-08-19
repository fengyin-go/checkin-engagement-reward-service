package handler

import (
	"net/http"
	"strconv"

	"checkin/pkg/httpx"
)

func (s *Server) registerLeaderboardRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/leaderboard/streak", s.streakLeaderboard)
	mux.HandleFunc("GET /api/leaderboard/total", s.totalDaysLeaderboard)
	mux.HandleFunc("GET /api/leaderboard/monthly", s.monthlyLeaderboard)
}

// parseLimit 解析 limit 参数，缺省返回 0（由 service 层取默认值）。
func parseLimit(r *http.Request) int {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	return limit
}

func (s *Server) streakLeaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := s.svc.StreakLeaderboard(parseLimit(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, entries)
}

func (s *Server) totalDaysLeaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := s.svc.TotalDaysLeaderboard(parseLimit(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, entries)
}

func (s *Server) monthlyLeaderboard(w http.ResponseWriter, r *http.Request) {
	entries, err := s.svc.MonthlyLeaderboard(r.URL.Query().Get("month"), parseLimit(r))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, entries)
}
