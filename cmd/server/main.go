package main

import (
	"github.com/gin-gonic/gin"

	httpDelivery "Baby-Rabbit/internal/delivery/http"
	"Baby-Rabbit/internal/pkg/logger"
	"Baby-Rabbit/internal/repository"
	"Baby-Rabbit/internal/service"
	"Baby-Rabbit/internal/usecase"
)

func main() {
	logger.Init()
	defer logger.Sync()

	logger.Log.Info("Starting Baby-Rabbit server...")

	manager := repository.NewQueueManager()
	useCase := usecase.NewQueueUseCase(manager)
	handler := httpDelivery.NewHandler(useCase)

	service.StartTTLCleaner(manager)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.POST("/queues", handler.CreateQueue)
	r.POST("/queues/:queue/push", handler.Push)
	r.POST("/queues/:queue/pop", handler.Pop)

	logger.Log.Info("HTTP server started on :8080")
	_ = r.Run(":8080")
}
