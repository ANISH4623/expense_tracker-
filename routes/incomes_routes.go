package routes

import (
	"awesomeProject1/controllers"
	"awesomeProject1/middleware"
	"github.com/gin-gonic/gin"
)

func IncomeRoutes(router *gin.Engine) {
	router.GET("/users/:userId/incomes", middleware.RequireAuth(), controllers.GetIncomes())
	router.POST("/users/:userId/incomes", middleware.RequireAuth(), controllers.CreateIncome())
	router.PUT("/users/:userId/incomes/:incomeId", middleware.RequireAuth(), controllers.UpdateIncome())
	router.GET("/users/:userId/incomes/:incomeId", middleware.RequireAuth(), controllers.GetIncome())
	router.DELETE("/users/:userId/incomes/:incomeId", middleware.RequireAuth(), controllers.DeleteIncome())
}
