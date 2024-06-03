package controllers

import (
	"awesomeProject1/database"
	"awesomeProject1/helpers"
	"awesomeProject1/middleware"
	"awesomeProject1/models"
	"awesomeProject1/serializers"
	"github.com/gin-gonic/gin"
	validate "github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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
		token, err := middleware.TokenController.CreateToken(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("jwt", token, 3600*24, "", "", false, true)
		c.JSON(http.StatusOK, gin.H{"data": UserResponse(user), "token": token})

	}
}
func User() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Request.Cookie("jwt")
		if err != nil {
			return
		}
		token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		}
		claims := token.Claims.(*jwt.StandardClaims)
		var user models.User
		database.Database.Db.Where("id = ?", claims.Subject).First(&user)
		c.JSON(200, user)
	}
}
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		helpers.ClearCookie(c, "jwt")
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
