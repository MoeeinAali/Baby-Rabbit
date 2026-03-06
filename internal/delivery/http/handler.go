package http

import (
	"Baby-Rabbit/internal/pkg/logger"
	"Baby-Rabbit/internal/usecase"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	usecase *usecase.QueueUseCase
}

func NewHandler(u *usecase.QueueUseCase) *Handler {
	return &Handler{usecase: u}
}

type CreateQueueReq struct {
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type CreateQueueResp struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateQueue(c *gin.Context) {
	var req CreateQueueReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	queueID, err := h.usecase.CreateQueue(req.Name, req.Capacity)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Infof("HTTP CreateQueue %s with ID %s", req.Name, queueID)
	c.JSON(200, CreateQueueResp{
		ID:   queueID,
		Name: req.Name,
	})
}

type PushReq struct {
	Value string `json:"value"`
	TTL   int    `json:"ttl"`
}

func (h *Handler) Push(c *gin.Context) {
	queueID := c.Param("queue")
	var req PushReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.usecase.Push(queueID, req.Value, req.TTL)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Infof("HTTP Push to queue %s", queueID)
	c.JSON(200, gin.H{"status": "ok"})
}

func (h *Handler) Pop(c *gin.Context) {
	queueID := c.Param("queue")

	msg, err := h.usecase.Pop(queueID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	logger.Log.Infof("HTTP Pop from queue %s, message %s", queueID, msg.ID)
	c.JSON(200, gin.H{"value": msg.Value})
}
