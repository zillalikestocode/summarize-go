package routes

import (
	"github.com/gin-gonic/gin"
	handler "github.com/zillalikestocode/summarize-api/internal/api/handlers"
)

func SetupRoutes(r *gin.Engine) {
	summaryHandler := handler.NewSummaryHandler()

	r.POST("/api/v1/summary", summaryHandler.Summarize)
}
