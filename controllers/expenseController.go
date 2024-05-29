package controllers

import (
	"awesomeProject1/database"
	"awesomeProject1/models"
	"awesomeProject1/serializers"
	"fmt"
	"github.com/gin-gonic/gin"
	validate "github.com/go-playground/validator/v10"
	"net/http"
)

func ExpenseResponse(expense models.Expense) serializers.ExpenseSerializer {
	return serializers.ExpenseSerializer{
		ID:       expense.ID,
		Amount:   expense.Amount,
		Category: expense.Category,
	}
}
func GetExpenses() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Initialize empty slice for expenses
		var expenses []models.Expense

		// Filter expenses based on authenticated user ID
		//result := database.Database.Db.Where("user_id = ?", userID).Find(&expenses)

		//// Handle database errors
		//if result.Error != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve expenses"})
		//	return
		//}
		user, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to retrieve user from context"})
			return
		}
		err := database.Database.Db.Model(&user).Association("Expenses").Find(&expenses)
		if err != nil {
			return
		}

		var Responsiveness []serializers.ExpenseSerializer
		for _, expense := range expenses {
			Responsiveness = append(Responsiveness, ExpenseResponse(expense))
		}
		// Return expenses in JSON response
		c.JSON(http.StatusOK, gin.H{"data": Responsiveness})
	}

}
func GetExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var expense models.Expense

		database.Database.Db.First(&expense, c.Param("expenseId"))
		if expense.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record Not Found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": ExpenseResponse(expense)})
	}
}
func CreateExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var expense models.Expense
		if err := c.Bind(&expense); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validate.New().Struct(expense); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to retrieve user from context"})
			return
		}
		tx := database.Database.Db.Begin()
		defer func() {
			if tx.Error != nil {
				tx.Rollback() // Rollback on error
			} else {
				tx.Commit() // Commit on success
			}
		}()
		expense.UserID = user.(models.User).ID
		if err := tx.Create(&expense).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create expense: %v", err)})
			return
		}
		c.JSON(201, gin.H{"data": expense})
	}
}

func UpdateExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var expense models.Expense
		if err := c.Bind(&expense); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validate.New().Struct(expense); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to retrieve user from context"})
			return
		}
		tx := database.Database.Db.Begin()
		defer func() {
			if tx.Error != nil {
				tx.Rollback() // Rollback on error
			} else {
				tx.Commit() // Commit on success
			}
		}()
		expense.UserID = user.(models.User).ID
		database.Database.Db.Save(models.Expense{
			ID:     expense.ID,
			UserID: expense.UserID, Amount: expense.Amount, Category: expense.Category,
		})
		c.JSON(http.StatusAccepted, gin.H{"data": expense})

	}
}
func DeleteExpense() gin.HandlerFunc {
	return func(c *gin.Context) {
		var expense models.Expense

		//user, ok := c.Get("user")
		//if !ok {
		//	c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to retrieve user from context"})
		//	return
		//}
		tx := database.Database.Db.Begin()
		defer func() {
			if tx.Error != nil {
				tx.Rollback() // Rollback on error
			} else {
				tx.Commit() // Commit on success
			}
		}()
		database.Database.Db.Delete(&expense, c.Param("expenseId"))
		c.JSON(200, gin.H{"message": "deleted successfully"})
	}
}
