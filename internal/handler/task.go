package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fashionparty/fast-task-manager-server/internal/middleware"
	"github.com/fashionparty/fast-task-manager-server/internal/model"
	"github.com/fashionparty/fast-task-manager-server/internal/service"
	"github.com/fashionparty/fast-task-manager-server/internal/validator"
)

type CreateTaskRequest struct {
	Title    string         `json:"title"`
	Content  string         `json:"content"`
	Priority model.Priority `json:"priority"`
	DueDate  *time.Time     `json:"duedate,omitempty"`
}

type UpdateTaskRequest struct {
	Title    *string         `json:"title,omitempty"`
	Content  *string         `json:"content,omitempty"`
	Priority *model.Priority `json:"priority,omitempty"`
	DueDate  *time.Time      `json:"duedate,omitempty"`
	IsDone   *bool           `json:"isdone"`
}

type TaskHandler interface {
	CreateTask(w http.ResponseWriter, r *http.Request)
	GetTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	ListTasks(w http.ResponseWriter, r *http.Request)
}

type taskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(service service.TaskService) TaskHandler {
	return &taskHandler{taskService: service}
}

func (h *taskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.ValidateRequired(req.Title, "title"); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	task, err := h.taskService.CreateTask(r.Context(), userID, req.Title, req.Content, req.Priority, req.DueDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, task)
}

func (h *taskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	taskID, err := parseID(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.taskService.GetTask(r.Context(), userID, taskID)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *taskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	taskID, err := parseID(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	task, err := h.taskService.UpdateTask(r.Context(), userID, taskID, req.Title, req.Content, req.Priority, req.IsDone, req.DueDate)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, task)
}

func (h *taskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	taskID, err := parseID(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	if err := h.taskService.DeleteTask(r.Context(), userID, taskID); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *taskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tasks, err := h.taskService.ListTasks(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, tasks)
}
