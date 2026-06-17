package controllers

import (
	"net/http"

	"backend/config"
	"backend/models"

	"github.com/gin-gonic/gin"
)

func CreatePermission(c *gin.Context) {

	var permission models.AccessPermission

	if err := c.ShouldBindJSON(&permission); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := config.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur création permission",
		})
		return
	}

	c.JSON(http.StatusCreated, permission)
}

func GetPermissions(c *gin.Context) {

	var permissions []models.AccessPermission

	config.DB.
		Preload("User").
		Preload("Door").
		Find(&permissions)

	c.JSON(http.StatusOK, permissions)
}
