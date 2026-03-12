package order

import (
	"encoding/json"
	"net/http"

	dto "github.com/barashF/lms/service-order/internal/handler/dto/order"
)

type Controller struct {
	service orderService
}

func NewController(s orderService) *Controller {
	return &Controller{service: s}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	id, err := c.service.Create(r.Context(), req.ToModel())
	if err != nil {
		c.writeError(w, err)
		return
	}

	c.writeJSON(w, http.StatusOK, map[string]any{
		"id":      id,
		"message": "Order created succesfily",
	})
}

func (c *Controller) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"Failed encode to response"}`, http.StatusInternalServerError)
	}
}

func (c *Controller) writeError(w http.ResponseWriter, err error) {
	msg, status := c.mapError(err)
	c.writeJSON(w, status, map[string]string{
		"error": msg,
	})
}

func (c *Controller) mapError(err error) (string, int) {
	switch err {
	default:
		return "Internal server error", http.StatusInternalServerError
	}
}
