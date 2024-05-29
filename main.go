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

func main() {

	helpers.LoadConfig(".env")
	database.Connect()
	router := gin.New()

	graphqlHandler := handler.New(&handler.Config{
		Schema: &schema,
		Pretty: true,
	})

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000/"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE"},
		AllowHeaders:     []string{"Orgin"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-lenght"},
		MaxAge:           12 * time.Hour,
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
