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

func IncomeResponse(Income models.Income) serializers.IncomeSerializer {
	return serializers.IncomeSerializer{
		ID:       Income.ID,
		Amount:   Income.Amount,
		Category: Income.Category,
	}
}
func GetIncomes() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Initialize empty slice for Incomes
		var Incomes []models.Income

		// Filter Incomes based on authenticated user ID
		//result := database.Database.Db.Where("user_id = ?", userID).Find(&Incomes)

		//// Handle database errors
		//if result.Error != nil {
		//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Incomes"})
		//	return
		//}
		user, ok := c.Get("user")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to retrieve user from context"})
			return
		}
		err := database.Database.Db.Model(&user).Association("Incomes").Find(&Incomes)
		if err != nil {
			return
		}

		var Responsiveness []serializers.IncomeSerializer
		for _, Income := range Incomes {
			Responsiveness = append(Responsiveness, IncomeResponse(Income))
		}
		// Return Incomes in JSON response
		c.JSON(http.StatusOK, gin.H{"data": Responsiveness})
	}

}
func GetIncome() gin.HandlerFunc {
	return func(c *gin.Context) {
		var Income models.Income

		database.Database.Db.First(&Income, c.Param("IncomeId"))
		if Income.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Record Not Found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": IncomeResponse(Income)})
	}
}
func CreateIncome() gin.HandlerFunc {
	return func(c *gin.Context) {
		var Income models.Income
		if err := c.Bind(&Income); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validate.New().Struct(Income); err != nil {
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
		Income.UserID = user.(models.User).ID
		if err := tx.Create(&Income).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create Income: %v", err)})
			return
		}
		c.JSON(201, gin.H{"data": Income})
	}
}

func UpdateIncome() gin.HandlerFunc {
	return func(c *gin.Context) {
		var Income models.Income
		if err := c.Bind(&Income); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validate.New().Struct(Income); err != nil {
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
		Income.UserID = user.(models.User).ID
		database.Database.Db.Save(models.Income{
			ID:     Income.ID,
			UserID: Income.UserID, Amount: Income.Amount, Category: Income.Category,
		})
		c.JSON(http.StatusAccepted, gin.H{"data": Income})

	}
}
func DeleteIncome() gin.HandlerFunc {
	return func(c *gin.Context) {
		var Income models.Income

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
		database.Database.Db.Delete(&Income, c.Param("IncomeId"))
		c.JSON(200, gin.H{"message": "deleted successfully"})
	}
}
