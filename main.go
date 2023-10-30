package main

import (
	"Reference/controller"
	"Reference/database"
	"Reference/util"

	"Reference/model"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnv() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func DBConnect() {
	database.ConnectDb()
	database.DB.AutoMigrate(&model.User{})
	// database.DB.AutoMigrate(&model.Role{})
	database.DB.AutoMigrate(&model.Library{})
	database.DB.AutoMigrate(&model.BookInventory{})
	database.DB.AutoMigrate(&model.RequestEvents{})
	database.DB.AutoMigrate(&model.IssueRegistry{})
	// seedData()
}

// func seedData() {
// 	var roles = []model.Role{{Name: "owner", Description: "Owner role"}, {Name: "admin", Description: "Administrator role"}, {Name: "reader", Description: "Reader role"}}
// 	// var user = []model.User{{Email: os.Getenv("ADMIN_EMAIL"), Password: os.Getenv("ADMIN_PASSWORD"), RoleID: 1}}
// 	database.DB.Save(&roles)
// 	// database.DB.Save(&user)
// }

func main() {
	router := gin.Default()

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000") // Adjust with your React app's address
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	loadEnv()
	DBConnect()

	authRoutes := router.Group("/auth/user")
	authRoutes.POST("/register", controller.Register)
	authRoutes.POST("/login", controller.Login)

	ownerRoutes := router.Group("/owner")
	ownerRoutes.Use(util.JWTAuth())
	ownerRoutes.POST("/library", controller.CreateLibrary)
	ownerRoutes.PUT("/library/:id", controller.UpdateLibrary)
	ownerRoutes.DELETE("/library/:id", controller.DeleteLibrary)

	adminRouter := router.Group("/admin")
	adminRouter.Use(util.JWTAuthAdmin())
	adminRouter.POST("/library/book", controller.AddBookInventory)
	adminRouter.PUT("/library/book/:isbn", controller.UpdateBook)
	adminRouter.DELETE("/library/book/:isbn", controller.DeleteBook)
	adminRouter.GET("/library/book/requests", controller.GetRequests)
	adminRouter.POST("/library/book/issue/approve/:request_id", controller.ApproveIssue)
	adminRouter.PUT("/library/book/issue/:issue_id", controller.ReturnIssue)
	adminRouter.GET("/library/book/issues", controller.GetAllIssues)

	openRouter := router.Group("/api")
	openRouter.Use(util.JWTAuthMiddleware())
	openRouter.GET("/library", controller.GetAllLibraries)
	openRouter.GET("/books/search", controller.SearchBooks)
	openRouter.GET("/library/books", controller.GetAllBook)
	openRouter.GET("/library/:library_id/books", controller.GetBookByLibraryID)
	openRouter.POST("/request/:book_isbn", controller.CreateRequest)
	openRouter.GET("/requests", controller.GetRequestsAsReader)

	// adminRoutes := router.Group("/admin")
	// adminRoutes.Use(util.JWTAuth())
	// adminRoutes.GET("/users", controller.GetUsers)
	// adminRoutes.GET("/user/:id", controller.GetUser)
	// adminRoutes.PUT("/user/:id", controller.UpdateUser)
	// adminRoutes.POST("/user/role", controller.CreateRole)
	// adminRoutes.GET("/user/roles", controller.GetRoles)
	// adminRoutes.PUT("/user/role/:id", controller.UpdateRole)

	router.Run(":8080")
}
