package handler

import (
	"net/http"

	"checkin/internal/model"
	"checkin/pkg/httpx"
)

func (s *Server) registerRewardRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/rewards", s.createReward)
	mux.HandleFunc("GET /api/rewards", s.listRewards)
	mux.HandleFunc("GET /api/rewards/{id}", s.getReward)
	mux.HandleFunc("PUT /api/rewards/{id}", s.updateReward)
	mux.HandleFunc("DELETE /api/rewards/{id}", s.deleteReward)
	mux.HandleFunc("GET /api/users/{id}/rewards", s.claimableRewards)
}

type rewardRequest struct {
	Name         string `json:"name"`
	Points       int    `json:"points"`
	RequiredDays int    `json:"required_days"`
	Description  string `json:"description"`
}

func (s *Server) createReward(w http.ResponseWriter, r *http.Request) {
	var req rewardRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	reward, err := s.svc.CreateReward(model.Reward{
		Name:         req.Name,
		Points:       req.Points,
		RequiredDays: req.RequiredDays,
		Description:  req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, reward)
}

func (s *Server) listRewards(w http.ResponseWriter, r *http.Request) {
	rewards, err := s.svc.ListRewards()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rewards)
}

func (s *Server) getReward(w http.ResponseWriter, r *http.Request) {
	reward, err := s.svc.GetReward(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, reward)
}

func (s *Server) updateReward(w http.ResponseWriter, r *http.Request) {
	var req rewardRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	reward, err := s.svc.UpdateReward(r.PathValue("id"), model.Reward{
		Name:         req.Name,
		Points:       req.Points,
		RequiredDays: req.RequiredDays,
		Description:  req.Description,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, reward)
}

func (s *Server) deleteReward(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteReward(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) claimableRewards(w http.ResponseWriter, r *http.Request) {
	rewards, err := s.svc.ClaimableRewards(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rewards)
}
