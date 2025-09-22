package handlers

import (
	"net/http"
	"notification/internal/models"
	"notification/internal/services"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	Service services.EmailService
}

func (h *EmailHandler) SendEmail(c *gin.Context) {
	var req models.EmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.Service.SendEmail(req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "email sent",
	})
}
