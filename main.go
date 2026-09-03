// curl.exe http://localhost:8080/... to call

// test createTask() logs
/*
$body = @{
    title = "Learn slog"
    description = "Structured logging practice"
    status = "pending"
} | ConvertTo-Json -Compress

Invoke-RestMethod `
-Method POST `
-Uri http://localhost:8080/tasks `
-ContentType "application/json" `
-Body $body

then check go run
*/

package main

import (
	"log/slog" // slog go import
	"os"       // for stdout

	"time"

	"net/http" // for http status codes
	"strconv"  // convert url id into integer

	"github.com/gin-gonic/gin" // gin framwork
)

// MARK: slog global var

var logger *slog.Logger // -> global logger variable

type Task struct { // -> the structure of the Task object
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateTaskInput struct { // update task inputs, no need ID for it since it just needs the contents
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}

// data stored in memory (TODO: add database later)
var tasks = []Task{
	{
		ID:          1,
		Title:       "Go Go Go",
		Description: "Basically Go",
		Status:      "pending",
	},
	{
		ID:          2,
		Title:       "Learning Gin",
		Description: "I love REST APIs",
		Status:      "pending",
	},
}

// simulating next task id
var nextID = 3

func main() { // -> main function for program

	// MARK: slog setup

	logger = setupLogger()

	logger.Info(
		"server starting",
	)

	router := gin.New() // -> creates the router without default loggers

	router.Use(gin.Recovery()) // recovery middleware

	router.Use(requestIDMiddleware()) // -> custom request ID middleware

	router.Use(slogMiddleware(logger)) // -> custom logger middleware

	router.GET("/health", healthHandler) // -> just to check if API is good (health check)

	// CRUD routes for tasks.
	router.GET("/tasks", getTasks)
	router.GET("/tasks/:id", getTask)
	router.POST("/tasks", createTask)
	router.PATCH("/tasks/:id", updateTask)
	router.DELETE("/tasks/:id", deleteTask)

	router.Run(":8080")
}

// healthHandler returns a simple response so I know the server is alive.
func healthHandler(c *gin.Context) { // -> health check handle

	reqID, _ := c.Get("request_id") // -> get the request id from the context

	// MARK: slog health check

	logger.Info( // logger for health check
		"health check requested",
		slog.String("endpoint", "/health"),
		slog.Any("request_id", reqID), // inject to the log
	)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func getTasks(c *gin.Context) { // -> get all tasks currently stored
	c.JSON(http.StatusOK, tasks)
}

func getTask(c *gin.Context) { // -> get a specific task from its id
	// c.Param("id") returns the id as a string.
	// Atoi converts string => int.
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{ // handling for incorrect task id
			"error": "Invalid task ID",
		})
		return
	}
	for _, task := range tasks { // loop until it finds the targeted task
		if task.ID == id {
			c.JSON(http.StatusOK, task)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{ // if loop finishes and no task is found
		"error": "Task not found",
	})
}

func createTask(c *gin.Context) { // to create a new task from JSON from the client
	var newTask Task // new task object

	// read request body and put values into newTask
	// then it gives gin the address of the variable so it can fill it in (&)
	err := c.ShouldBindJSON(&newTask)

	if err != nil {

		// MARK: slog warn

		logger.Warn( // using slog for errors
			"invalid task payload",
			slog.Any("error", err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// applying id logic
	newTask.ID = nextID
	nextID++

	// add new task in the memory storage
	tasks = append(tasks, newTask)

	// MARK: slog info

	// logging
	logger.Info(
		"task created",
		slog.Int("task_id", newTask.ID),
		slog.String("title", newTask.Title),
		slog.String("status", newTask.Status),
	)

	// new resource was successfully created (201 status code)
	c.JSON(http.StatusCreated, newTask)
}

func updateTask(c *gin.Context) { // to update the fields on existing tasks
	id, err := strconv.Atoi(c.Param("id")) // same logic as getting the task

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	var input UpdateTaskInput // create a new input object to store input values

	// put the JSON body into our update input struct
	err = c.ShouldBindJSON(&input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// using i as the task index, update based on the input values
	for i := range tasks {
		if tasks[i].ID == id {

			if input.Title != nil { // -> nil means that the field is not included in the request
				tasks[i].Title = *input.Title // -> pointer is used to reference the actual value (string)
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

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Task not found",
	})
}

func deleteTask(c *gin.Context) { // delete a task based on id
	id, err := strconv.Atoi(c.Param("id")) // same logic as get task

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid task ID",
		})
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...) // keeping everything before and after i, then join them together

			c.Status(http.StatusNoContent) // 204 -> delete works
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{
		"error": "Task not found",
	})
}

func setupLogger() *slog.Logger {

	handler := slog.NewJSONHandler(
		os.Stdout,
		nil,
	)

	logger := slog.New(handler)

	return logger.With(
		slog.String("service", "task-api"),
		slog.String("environment", "development"),
	)

}

// MARK: slog middleware

func slogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		reqID, _ := c.Get("request_id") // -> get the request id from the context

		logger.Info(
			"HTTP request",
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("latency", latency.String()),
			slog.Any("request_id", reqID), // inject to the log
		)
	}
}

// MARK: context logging

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. generate a mock Request ID (in production, use UUIDs)
		reqID := "req-" + strconv.FormatInt(time.Now().UnixNano(), 10)

		// 2. add to Gin's internal context so handlers can access it
		c.Set("request_id", reqID)

		// 3. add it to the HTTP response header (good practice)
		c.Writer.Header().Set("X-Request-ID", reqID)

		c.Next()
	}
}
