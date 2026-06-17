package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
	}

	users := api.Group("/users")
	users.Use(middlewares.AuthMiddleware())
	{
		users.GET("/profile", controllers.Profile)
	}

	doors := api.Group("/doors")
	doors.Use(middlewares.AuthMiddleware())
	{
		doors.POST("/", controllers.CreateDoor)
		doors.GET("/", controllers.GetDoors)

		doors.POST("/:id/unlock", controllers.UnlockDoor)
		doors.POST("/:id/lock", controllers.LockDoor)
	}

	permissions := api.Group("/permissions")
	permissions.Use(middlewares.AuthMiddleware())
	{
		permissions.POST("/", controllers.CreatePermission)
		permissions.GET("/", controllers.GetPermissions)
	}

	access := api.Group("/access")
	access.Use(middlewares.AuthMiddleware())
	{
		access.POST("/open", controllers.OpenAccess)
		access.GET("/logs", controllers.GetAccessLogs)
	}

	faces := api.Group("/faces")
	faces.Use(middlewares.AuthMiddleware())
	{
		faces.POST("/register", controllers.RegisterFace)
		faces.GET("/", controllers.GetFaces)
		faces.POST("/verify", controllers.VerifyFace)
		faces.POST("/verify-image", controllers.VerifyFaceImage)
	}

	admin := api.Group("/admin")
	admin.Use(middlewares.AuthMiddleware())
	{
		admin.POST("/retention/clean-face-images", controllers.CleanExpiredFaceImages)
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
