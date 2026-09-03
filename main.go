package main

import (
	"go-learn-it/controllers/healthcheckcontroller"
	"go-learn-it/controllers/taskcontroller"
	"go-learn-it/middlewares"
	"go-learn-it/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Initialize utilities (Logger)
	logger := utils.SetupLogger()
	logger.Info("server starting")

	// 2. Initialize Controllers with Dependency Injection
	healthCtrl := healthcheckcontroller.NewHealthCheckController(logger)
	taskCtrl := taskcontroller.NewTaskController(logger)

	// 3. Setup Router & Middlewares
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middlewares.RequestIDMiddleware())
	router.Use(middlewares.SlogMiddleware(logger))

	// 4. Mount Routes
	router.GET("/health", healthCtrl.HealthCheck)
	router.GET("/tasks", taskCtrl.GetTasks)
	router.GET("/tasks/:id", taskCtrl.GetTask)
	router.POST("/tasks", taskCtrl.CreateTask)
	router.PATCH("/tasks/:id", taskCtrl.UpdateTask)
	router.DELETE("/tasks/:id", taskCtrl.DeleteTask)

	// 5. Start Server
	router.Run(":8080")
}
