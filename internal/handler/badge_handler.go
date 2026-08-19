package handler

import (
	"net/http"

	"checkin/internal/model"
	"checkin/pkg/httpx"
)

func (s *Server) registerBadgeRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/badges", s.createBadge)
	mux.HandleFunc("GET /api/badges", s.listBadges)
	mux.HandleFunc("GET /api/badges/{id}", s.getBadge)
	mux.HandleFunc("PUT /api/badges/{id}", s.updateBadge)
	mux.HandleFunc("DELETE /api/badges/{id}", s.deleteBadge)
	mux.HandleFunc("GET /api/users/{id}/badges", s.userBadges)
}

type badgeRequest struct {
	Name        string          `json:"name"`
	Type        model.BadgeType `json:"type"`
	Threshold   int             `json:"threshold"`
	Description string          `json:"description"`
}

func (s *Server) createBadge(w http.ResponseWriter, r *http.Request) {
	var req badgeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.CreateBadge(model.Badge{
		Name:        req.Name,
		Type:        req.Type,
		Threshold:   req.Threshold,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, b)
}

func (s *Server) listBadges(w http.ResponseWriter, r *http.Request) {
	badges, err := s.svc.ListBadges()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, badges)
}

func (s *Server) getBadge(w http.ResponseWriter, r *http.Request) {
	b, err := s.svc.GetBadge(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) updateBadge(w http.ResponseWriter, r *http.Request) {
	var req badgeRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	b, err := s.svc.UpdateBadge(r.PathValue("id"), model.Badge{
		Name:        req.Name,
		Type:        req.Type,
		Threshold:   req.Threshold,
		Description: req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, b)
}

func (s *Server) deleteBadge(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteBadge(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) userBadges(w http.ResponseWriter, r *http.Request) {
	progress, err := s.svc.UserBadges(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, progress)
}
