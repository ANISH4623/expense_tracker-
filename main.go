package main

import (
	"awesomeProject1/database"
	"awesomeProject1/helpers"
	"awesomeProject1/middleware"
	"awesomeProject1/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/handler"
	"os"
	"time"
)

// Defining the Playground handler

func main() {
	graphqlHandler := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	helpers.LoadConfig(".env")
	database.Connect()
	router := gin.New()
	middleware.TokenController = helpers.NewJWTToken(helpers.AppConfig.SECRET_KEY)

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Access-Control-Allow-Headers", "Content-Type", "Content-Length", " Accept-Encoding", " X-CSRF-Token", "Authorization", "accept", " origin", "Cache-Control", " X-Requested-With", "set-cookie"},
		AllowWildcard:    false,
		MaxAge:           12 * time.Hour,
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:3000"},
	}))

	port := os.Getenv("PORT")
	routes.AuthRoutes(router)
	routes.Expense_Routes(router)
	routes.IncomeRoutes(router)
	router.POST("/graphql", middleware.RequireAuth(), func(c *gin.Context) {
		graphqlHandler.ServeHTTP(c.Writer, c.Request)
	})
	err := router.Run(":" + port)
	if err != nil {
		return
	}
}
