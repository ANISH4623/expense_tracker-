package routes

import (
	"awesomeProject1/controllers"
	"awesomeProject1/middleware"
	"github.com/gin-gonic/gin"
)

func Expense_Routes(router *gin.Engine) {
	router.GET("/users/:userId/expenses", middleware.RequireAuth(), controllers.GetExpenses())
	router.POST("/users/:userId/expenses", middleware.RequireAuth(), controllers.CreateExpense())
	router.PUT("/users/:userId/expenses/:expenseId", middleware.RequireAuth(), controllers.UpdateExpense())
	router.GET("/users/:userId/expenses/:expenseId", middleware.RequireAuth(), controllers.GetExpense())
	router.DELETE("/users/:userId/expenses/:expenseId", middleware.RequireAuth(), controllers.DeleteExpense())
}
