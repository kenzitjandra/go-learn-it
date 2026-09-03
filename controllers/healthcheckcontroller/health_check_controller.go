package healthcheckcontroller

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthCheckController struct {
	Logger *slog.Logger
}

func NewHealthCheckController(logger *slog.Logger) *HealthCheckController {
	return &HealthCheckController{Logger: logger}
}

func (ctrl *HealthCheckController) HealthCheck(c *gin.Context) {
	reqID, _ := c.Get("request_id")

	ctrl.Logger.Info(
		"health check requested",
		slog.String("endpoint", "/health"),
		slog.Any("request_id", reqID),
	)

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
