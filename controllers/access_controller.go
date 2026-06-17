package controllers

import (
	"net/http"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

type OpenAccessInput struct {
	DoorID uint `json:"door_id"`
}

func OpenAccess(c *gin.Context) {
	var input OpenAccessInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Utilisateur non authentifié",
		})
		return
	}

	userID := userIDValue.(uint)

	var permission models.AccessPermission

	err := config.DB.
		Where("user_id = ? AND door_id = ? AND can_access = ?", userID, input.DoorID, true).
		First(&permission).Error

	if err != nil {
		log := models.AccessLog{
			UserID:   &userID,
			DoorID:   input.DoorID,
			Status:   "DENIED",
			Reason:   "Permission inexistante ou refusée",
			SourceIP: c.ClientIP(),
		}
		config.DB.Create(&log)

		c.JSON(http.StatusForbidden, gin.H{
			"access":  "DENIED",
			"message": "Accès refusé",
		})
		return
	}

	var door models.Door

	if err := config.DB.First(&door, input.DoorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Porte introuvable"})
		return
	}

	door.IsLocked = false
	config.DB.Save(&door)

	log := models.AccessLog{
		UserID:   &userID,
		DoorID:   input.DoorID,
		Status:   "GRANTED",
		Reason:   "Permission valide",
		SourceIP: c.ClientIP(),
	}
	config.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"access":  "GRANTED",
		"message": "Accès autorisé, porte déverrouillée",
		"door":    door,
	})
}

func GetAccessLogs(c *gin.Context) {
	var logs []models.AccessLog

	config.DB.
		Preload("User").
		Preload("Door").
		Order("created_at DESC").
		Find(&logs)

	c.JSON(http.StatusOK, logs)
}
