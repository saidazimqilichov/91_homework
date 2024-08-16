package handlers

import (
	"cicd/models"
	"cicd/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.repo.CreateTask(c.Request.Context(), &task); err != nil {
		c.JSON(500, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id := c.Param("id")

	task, err := h.repo.GetTaskByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.repo.GetTasks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var update map[string]interface{}
	if err := c.ShouldBindBodyWithJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := h.repo.UpdateTask(c.Request.Context(), id, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteTask(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
