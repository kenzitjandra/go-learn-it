package taskcontroller

import (
	"log/slog"
	"net/http"
	"strconv"

	"go-learn-it/models"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	Logger *slog.Logger
}

func NewTaskController(logger *slog.Logger) *TaskController {
	return &TaskController{Logger: logger}
}

// Temporary in-memory storage (will move to repositories later)
var tasks = []models.Task{
	{ID: 1, Title: "Go Go Go", Description: "Basically Go", Status: "pending"},
	{ID: 2, Title: "Learning Gin", Description: "I love REST APIs", Status: "pending"},
}
var nextID = 3

func (ctrl *TaskController) GetTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}

func (ctrl *TaskController) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	for _, task := range tasks {
		if task.ID == id {
			c.JSON(http.StatusOK, task)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func (ctrl *TaskController) CreateTask(c *gin.Context) {
	reqID, _ := c.Get("request_id")
	var newTask models.Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		ctrl.Logger.Warn(
			"invalid task payload",
			slog.Any("error", err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask.ID = nextID
	nextID++
	tasks = append(tasks, newTask)

	ctrl.Logger.Info(
		"task created",
		slog.Int("task_id", newTask.ID),
		slog.String("title", newTask.Title),
		slog.String("status", newTask.Status),
		slog.Any("request_id", reqID),
	)

	c.JSON(http.StatusCreated, newTask)
}

func (ctrl *TaskController) UpdateTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var input models.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			if input.Title != nil {
				tasks[i].Title = *input.Title
			}
			if input.Description != nil {
				tasks[i].Description = *input.Description
			}
			if input.Status != nil {
				tasks[i].Status = *input.Status
			}
			c.JSON(http.StatusOK, tasks[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}

func (ctrl *TaskController) DeleteTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
}
