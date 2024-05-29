package controllers

import (
	"awesomeProject1/database"
	"awesomeProject1/helpers"
	"awesomeProject1/models"
	"awesomeProject1/serializers"
	"github.com/gin-gonic/gin"
	validate "github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func UserResponse(user models.User) serializers.UserSerializer {
	return serializers.UserSerializer{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.String(),
		UpdatedAt: user.UpdatedAt.String(),
	}
}
func Signup() gin.HandlerFunc {
	return func(context *gin.Context) {
		var existingUser models.User
		var user models.User
		if err := context.Bind(&user); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := validate.New().Struct(user); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		database.Database.Db.Find(&existingUser, "email = ?", user.Email)
		if existingUser.ID != 0 {
			context.JSON(http.StatusBadRequest, gin.H{"Error": "User Already Exists "})
			return
		}
		password, _ := hashPassword(user.Password)
		user.Password = password
		database.Database.Db.Create(&user)
		var ResponseUser = UserResponse(user)
		context.JSON(http.StatusCreated, ResponseUser)
	}

}
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var data map[string]string
		var user models.User
		if err := c.Bind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		database.Database.Db.Where("email = ?", data["email"]).First(&user)
		if user.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No user found Please Signup"})
			return
		}
		if match := checkPasswordHash(data["password"], user.Password); !match {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Credentials"})
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})
		tokenString, err := token.SignedString([]byte(helpers.AppConfig.SECRET_KEY))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create token",
			})
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
		c.JSON(
			http.StatusOK,
			gin.H{"data": UserResponse(user)},
		)

	}
}
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		helpers.ClearCookie(c, "authorization")
		c.JSON(http.StatusOK, gin.H{"Message": "Logged out Successfully"})
	}
}
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
