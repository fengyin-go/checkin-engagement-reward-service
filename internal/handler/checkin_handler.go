package handler

import (
	"net/http"

	"checkin/pkg/httpx"
)

func (s *Server) registerCheckinRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/checkins", s.checkIn)
	mux.HandleFunc("DELETE /api/checkins/{id}", s.deleteCheckin)
	mux.HandleFunc("GET /api/users/{id}/summary", s.checkinSummary)
	mux.HandleFunc("GET /api/users/{id}/checkins", s.listCheckins)
	mux.HandleFunc("GET /api/stats/checkins", s.monthlyStats)
}

type checkInRequest struct {
	UserID string `json:"user_id"`
}

func (s *Server) checkIn(w http.ResponseWriter, r *http.Request) {
	var req checkInRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CheckIn(req.UserID)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) deleteCheckin(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteCheckin(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) checkinSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := s.svc.GetSummary(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, summary)
}

func (s *Server) listCheckins(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	items, total, err := s.svc.ListCheckins(r.PathValue("id"), r.URL.Query().Get("month"), pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) monthlyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := s.svc.MonthlyStats(r.URL.Query().Get("month"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, stats)
}
