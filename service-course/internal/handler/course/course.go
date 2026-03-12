package course

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/barashF/lms/service-course/internal/handler/course/validation"
	dto "github.com/barashF/lms/service-course/internal/handler/dto/course"
	"github.com/barashF/lms/service-course/internal/model"
)

type Controller struct {
	service courseService
}

func NewController(s courseService) *Controller {
	return &Controller{service: s}
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	courseID, err := uuid.Parse(idStr)
	if err != nil {
		c.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid course id format",
		})
	}

	course, err := c.service.FetchByID(r.Context(), courseID)
	if err != nil {
		c.writeError(w, err)
	}

	c.writeJSON(w, http.StatusOK, dto.ModelToResponse(*course))
}

func (c *Controller) GetMany(w http.ResponseWriter, r *http.Request) {
	courses, err := c.service.FetchAll(r.Context())
	if err != nil {
		c.writeError(w, err)
	}

	responseCourses := make([]dto.Course, len(courses))
	for i, c := range courses {
		responseCourses[i] = dto.ModelToResponse(*c)
	}

	c.writeJSON(w, http.StatusOK, responseCourses)
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if err := validation.ValidateRequest(&req); err != nil {
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

	c.writeJSON(w, http.StatusCreated, map[string]any{
		"id":      id,
		"message": "Course created succesfilly",
	})
}

func (c *Controller) writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

func (c *Controller) writeError(w http.ResponseWriter, err error) {
	message, status := c.mapError(err)
	c.writeJSON(w, status, map[string]string{"error": message})
}

func (c *Controller) mapError(err error) (string, int) {
	switch {
	case errors.Is(err, model.ErrNoCourseFound):
		return "Course not found", http.StatusNotFound
	default:
		return "Internal server error", http.StatusInternalServerError
	}
}
