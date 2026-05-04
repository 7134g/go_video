package api

import (
	"net/http"
	"strconv"

	"go_video/internal/model"
	"go_video/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	svc *service.TaskService
}

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{svc: service.NewTaskService()}
}

type IDReq struct {
	ID uint `json:"id"`
}

type CreateTaskReq struct {
	Name   string `json:"name" binding:"required"`
	URL    string `json:"url" binding:"required"`
	Header string `json:"header"`
	Type   string `json:"type" binding:"required"`
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &model.Task{
		Name:   req.Name,
		URL:    req.URL,
		Header: req.Header,
		Type:   req.Type,
	}
	if err := h.svc.Create(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	var req IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.Delete(req.ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

type UpdateTaskReq struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	URL    string `json:"url"`
	Header string `json:"header"`
	Type   string `json:"type"`
}

func (h *TaskHandler) Update(c *gin.Context) {
	var req UpdateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.svc.GetByID(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.URL != "" {
		task.URL = req.URL
	}
	task.Header = req.Header
	if req.Type != "" {
		task.Type = req.Type
	}

	if err := h.svc.Update(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Start(c *gin.Context) {
	count, err := h.svc.StartTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"started": count})
}

func (h *TaskHandler) Pause(c *gin.Context) {
	var req IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.PauseTask(req.ID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "paused"})
}

func (h *TaskHandler) Retry(c *gin.Context) {
	var req IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.RetryTask(req.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "retrying"})
}

func (h *TaskHandler) StartOne(c *gin.Context) {
	var req IDReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.StartTask(req.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "started"})
}

func (h *TaskHandler) List(c *gin.Context) {
	statusStr := c.Query("status")
	var tasks []model.Task
	var err error

	if statusStr != "" {
		status, _ := strconv.Atoi(statusStr)
		tasks, err = h.svc.GetByStatus(model.TaskStatus(status))
	} else {
		tasks, err = h.svc.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
