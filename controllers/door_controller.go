package controllers

import (
	"net/http"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func CreateDoor(c *gin.Context) {

	var door models.Door

	if err := c.ShouldBindJSON(&door); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := config.DB.Create(&door).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur création porte",
		})
		return
	}

	c.JSON(http.StatusCreated, door)
}

func GetDoors(c *gin.Context) {

	var doors []models.Door

	config.DB.Find(&doors)

	c.JSON(http.StatusOK, doors)
}

func UnlockDoor(c *gin.Context) {

	id := c.Param("id")

	var door models.Door

	if err := config.DB.First(&door, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Porte introuvable",
		})
		return
	}

	door.IsLocked = false

	config.DB.Save(&door)

	c.JSON(200, gin.H{
		"message": "Porte déverrouillée",
		"door":    door,
	})
}

func LockDoor(c *gin.Context) {

	id := c.Param("id")

	var door models.Door

	if err := config.DB.First(&door, id).Error; err != nil {
		c.JSON(404, gin.H{
			"error": "Porte introuvable",
		})
		return
	}

	door.IsLocked = true

	config.DB.Save(&door)

	c.JSON(200, gin.H{
		"message": "Porte verrouillée",
		"door":    door,
	})
}
