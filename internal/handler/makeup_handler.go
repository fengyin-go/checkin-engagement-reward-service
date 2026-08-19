package handler

import (
	"net/http"

	"checkin/pkg/httpx"
)

func (s *Server) registerMakeupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/makeup-cards", s.grantMakeupCard)
	mux.HandleFunc("GET /api/users/{id}/makeup-card", s.getMakeupCard)
	mux.HandleFunc("POST /api/makeups", s.makeupCheckIn)
	mux.HandleFunc("GET /api/users/{id}/makeups", s.listMakeups)
}

type grantMakeupCardRequest struct {
	UserID string `json:"user_id"`
	Amount int    `json:"amount"`
}

func (s *Server) grantMakeupCard(w http.ResponseWriter, r *http.Request) {
	var req grantMakeupCardRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	card, err := s.svc.GrantMakeupCard(req.UserID, req.Amount)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, card)
}

func (s *Server) getMakeupCard(w http.ResponseWriter, r *http.Request) {
	card, err := s.svc.GetMakeupCard(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, card)
}

type makeupCheckInRequest struct {
	UserID string `json:"user_id"`
	Date   string `json:"date"`
}

func (s *Server) makeupCheckIn(w http.ResponseWriter, r *http.Request) {
	var req makeupCheckInRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.MakeupCheckIn(req.UserID, req.Date)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listMakeups(w http.ResponseWriter, r *http.Request) {
	items, err := s.svc.ListMakeups(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, items)
}
