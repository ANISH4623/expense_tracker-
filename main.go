package main

import (
	"awesomeProject1/database"
	"awesomeProject1/helpers"
	"awesomeProject1/routes"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {

	helpers.LoadConfig(".env")
	database.Connect()
	router := gin.New()
	port := os.Getenv("PORT")
	routes.AuthRoutes(router)
	routes.Expense_Routes(router)
	routes.IncomeRoutes(router)
	err := router.Run(":" + port)
	if err != nil {
		return
	}
}
