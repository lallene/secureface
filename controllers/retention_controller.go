package controllers

import (
	"net/http"
	"os"
	"time"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func CleanExpiredFaceImages(c *gin.Context) {
	var profiles []models.FaceProfile

	now := time.Now()

	config.DB.
		Where("image_path <> '' AND image_expires_at <= ? AND image_deleted_at IS NULL", now).
		Find(&profiles)

	deletedCount := 0

	for _, profile := range profiles {
		if err := os.Remove(profile.ImagePath); err == nil {
			deletedAt := time.Now()

			profile.ImagePath = ""
			profile.ImageDeletedAt = &deletedAt

			config.DB.Save(&profile)

			deletedCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Nettoyage terminé",
		"deleted_count": deletedCount,
	})
}
