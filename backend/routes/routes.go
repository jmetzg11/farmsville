package routes

import (
	"farmsville/backend/api"
	"farmsville/backend/database"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupAPIRoutes(router *gin.Engine) {
	setupCORS(router)

	handler := api.NewHandler(database.DB)

	apiRouter := router.Group("/api")
	{
		apiRouter.GET("/hello", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Hello from Go!"})
		})
		apiRouter.GET("/show_auth", handler.ShowAuth)

		// customers
		apiRouter.GET("/items", handler.GetItems)
		apiRouter.GET("/blogs", handler.GetBlogs)

		// auth
		apiRouter.POST("/auth", handler.SendAuth)
		apiRouter.POST("/auth/verify", handler.VerifyAuth)
		apiRouter.GET("/auth/me", handler.AuthMe)
		apiRouter.POST("/auth/login", handler.LoginWithPassword)
		apiRouter.POST("/auth/create", handler.CreateAccount)
		apiRouter.POST("/auth/code-to-reset-password", handler.SendCodeToResetPassword)
		apiRouter.POST("/auth/reset-password", handler.ResetPassword)
		apiRouter.GET("/auth/logout", handler.Logout)

		// messages
		apiRouter.GET("/messages", handler.GetMessages)

		authRoutes := apiRouter.Group("/")
		authRoutes.Use(handler.AuthMiddleware())
		{
			// customers
			authRoutes.POST("/items/claim", handler.MakeClaim)
		}

		// admin
		adminRoutes := apiRouter.Group("/")
		adminRoutes.Use(handler.AdminMiddleware())
		{
			// inventory
			adminRoutes.POST("/items/update", handler.UpdateItem)
			adminRoutes.POST("/items/remove", handler.RemoveItem)
			adminRoutes.POST("/items/create", handler.CreateItem)
			adminRoutes.POST("/items/admin-claim", handler.AdminClaimItem)
			adminRoutes.POST("/claimed-item/remove", handler.RemoveClaimedItem)

			// users
			adminRoutes.GET("/users", handler.GetUsers)
			adminRoutes.POST("/users/update", handler.UpdateUser)
			adminRoutes.POST("/users/remove", handler.RemoveUser)
			adminRoutes.POST("/users/create", handler.CreateUser)

			// messages
			adminRoutes.POST("/post-message", handler.PostMessage)
			adminRoutes.DELETE("/messages/:id", handler.DeleteMessage)
			adminRoutes.POST("/send-email", handler.SendEmail)

			// blog
			adminRoutes.POST("/post-blog", handler.PostBlog)
			adminRoutes.GET("/get-blog-titles", handler.GetBlogTitles)
			adminRoutes.GET("/get-blog/:id", handler.GetBlogById)
			adminRoutes.POST("/edit-blog", handler.EditBlog)
		}
	}
}

func setupCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow local frontend in dev
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func SetupStaticRoutes(router *gin.Engine) {
	router.Use(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Next() // Let API routes handle it
			return
		}

		// Try to serve static files
		filePath := "./frontend/build" + c.Request.URL.Path
		if _, err := os.Stat(filePath); err == nil {
			c.File(filePath)
			c.Abort()
			return
		}
		// If not a file, serve index.html (for SvelteKit routing)
		c.File("./frontend/build/index.html")
		c.Abort()
	})
}
