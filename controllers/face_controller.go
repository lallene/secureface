package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"backend/config"
	"backend/models"
	"backend/utils"

	"github.com/gin-gonic/gin"
)

type FaceRegisterInput struct {
	UserID     uint   `json:"user_id"`
	FaceVector string `json:"face_vector"`
	ImagePath  string `json:"image_path"`
}

type FaceAIResponse struct {
	Match      bool    `json:"match"`
	UserID     uint    `json:"user_id"`
	Confidence float64 `json:"confidence"`
	Message    string  `json:"message"`
}

type FaceVerifyInput struct {
	UserID uint `json:"user_id"`
	DoorID uint `json:"door_id"`
}

func RegisterFace(c *gin.Context) {
	var input FaceRegisterInput
	expiresAt := time.Now().AddDate(0, 0, 7)

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	encryptedVector, err := utils.EncryptFaceVector(input.FaceVector)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Erreur chiffrement empreinte faciale",
		})
		return
	}

	face := models.FaceProfile{
		UserID:         input.UserID,
		FaceVector:     encryptedVector,
		ImagePath:      "",
		ImageExpiresAt: &expiresAt,
		IsVerified:     true,
	}

	if err := config.DB.Create(&face).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création profil facial"})
		return
	}

	c.JSON(http.StatusCreated, face)
}

func GetFaces(c *gin.Context) {
	var faces []models.FaceProfile

	config.DB.
		Preload("User").
		Find(&faces)

	c.JSON(http.StatusOK, faces)
}

func VerifyFace(c *gin.Context) {
	var input FaceVerifyInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	var face models.FaceProfile

	if err := config.DB.
		Where("user_id = ? AND is_verified = ?", input.UserID, true).
		First(&face).Error; err != nil {

		log := models.AccessLog{
			UserID:   &input.UserID,
			DoorID:   input.DoorID,
			Status:   "DENIED",
			Reason:   "Profil facial introuvable ou non vérifié",
			SourceIP: c.ClientIP(),
		}
		config.DB.Create(&log)

		c.JSON(http.StatusForbidden, gin.H{
			"access":  "DENIED",
			"message": "Visage non reconnu ou non vérifié",
		})
		return
	}

	var permission models.AccessPermission

	if err := config.DB.
		Where("user_id = ? AND door_id = ? AND can_access = ?", input.UserID, input.DoorID, true).
		First(&permission).Error; err != nil {

		log := models.AccessLog{
			UserID:   &input.UserID,
			DoorID:   input.DoorID,
			Status:   "DENIED",
			Reason:   "Visage reconnu mais permission refusée",
			SourceIP: c.ClientIP(),
		}
		config.DB.Create(&log)

		c.JSON(http.StatusForbidden, gin.H{
			"access":  "DENIED",
			"message": "Visage reconnu mais accès refusé",
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
		UserID:   &input.UserID,
		DoorID:   input.DoorID,
		Status:   "GRANTED",
		Reason:   "Visage reconnu et permission valide",
		SourceIP: c.ClientIP(),
	}
	config.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"access":  "GRANTED",
		"message": "Visage reconnu, accès autorisé",
		"door":    door,
	})
}

func VerifyFaceImage(c *gin.Context) {
	doorIDStr := c.PostForm("door_id")

	doorID64, err := strconv.ParseUint(doorIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "door_id invalide"})
		return
	}
	doorID := uint(doorID64)

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image obligatoire"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lecture image"})
		return
	}
	defer openedFile.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("image", file.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur préparation image"})
		return
	}

	if _, err := io.Copy(part, openedFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur copie image"})
		return
	}

	writer.Close()

	faceAIURL := os.Getenv("FACE_AI_URL")
	if faceAIURL == "" {
		faceAIURL = "http://localhost:8000/recognize"
	}

	req, err := http.NewRequest("POST", faceAIURL, &requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur création requête IA"})
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service IA indisponible"})
		return
	}
	defer resp.Body.Close()

	var aiResponse FaceAIResponse

	if err := json.NewDecoder(resp.Body).Decode(&aiResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Réponse IA invalide"})
		return
	}

	if !aiResponse.Match {
		confidence := aiResponse.Confidence
		log := models.AccessLog{
			UserID:          nil,
			DoorID:          doorID,
			Status:          "DENIED",
			Reason:          "Visage non reconnu par IA",
			SourceIP:        c.ClientIP(),
			ConfidenceScore: &confidence,
		}
		config.DB.Create(&log)

		c.JSON(http.StatusForbidden, gin.H{
			"access":     "DENIED",
			"message":    "Visage non reconnu",
			"confidence": aiResponse.Confidence,
		})
		return
	}

	var permission models.AccessPermission

	if err := config.DB.
		Where("user_id = ? AND door_id = ? AND can_access = ?", aiResponse.UserID, doorID, true).
		First(&permission).Error; err != nil {

		confidence := aiResponse.Confidence

		log := models.AccessLog{
			UserID:          &aiResponse.UserID,
			DoorID:          doorID,
			Status:          "DENIED",
			Reason:          "Visage reconnu mais permission refusée",
			SourceIP:        c.ClientIP(),
			ConfidenceScore: &confidence,
		}
		config.DB.Create(&log)

		c.JSON(http.StatusForbidden, gin.H{
			"access":     "DENIED",
			"message":    "Visage reconnu mais accès refusé",
			"user_id":    aiResponse.UserID,
			"confidence": aiResponse.Confidence,
		})
		return
	}

	var door models.Door

	if err := config.DB.First(&door, doorID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Porte introuvable"})
		return
	}

	door.IsLocked = false
	config.DB.Save(&door)

	confidence := aiResponse.Confidence

	log := models.AccessLog{
		UserID:          &aiResponse.UserID,
		DoorID:          doorID,
		Status:          "GRANTED",
		Reason:          "Visage reconnu par IA et permission valide",
		SourceIP:        c.ClientIP(),
		ConfidenceScore: &confidence,
	}
	config.DB.Create(&log)

	c.JSON(http.StatusOK, gin.H{
		"access":     "GRANTED",
		"message":    "Visage reconnu, accès autorisé",
		"user_id":    aiResponse.UserID,
		"confidence": aiResponse.Confidence,
		"door":       door,
	})
}
