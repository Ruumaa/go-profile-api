package controllers

import (
	"errors"
	"go-web/database"
	"go-web/models"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetAllUsers(c *gin.Context) {
	var users []models.User

	database.DB.Preload("Photos").Find(&users)

	c.JSON(http.StatusOK, gin.H{"message": "GET users success", "data": users})
}

func Register(c *gin.Context) {
	var body models.User

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Username, email, and password are required", "error": err.Error()})
		return
	}

	// Menangani len password <6
	if len(body.Password) < 6 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password must be at least 6 characters"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "failed to hashed password",
			"error":   err.Error(),
		})
		return
	}

	hashed := models.User{Username: body.Username, Email: body.Email, Password: string(hash)}

	if result := database.DB.Create(&hashed); result.Error != nil {
		// Cek apakah error disebabkan oleh uniqueness constraint pada kolom email
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Email is already taken"})
			return
		}

		// Tangani error lainnya sesuai kebutuhan
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Create new user success", "data": hashed})
}

func Login(c *gin.Context) {
	var body struct {
		Username string
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	var user models.User
	database.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid email or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to create token", "error": err.Error()})
		return
	}

	// add token
	user.Token = tokenString
	if result := database.DB.Save(&user); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to add token", "error": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login success"})
}

func EditUser(c *gin.Context) {
	var body models.User
	id := c.Param("id")

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if len(body.Password) < 6 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Password must be at least 6 characters"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "failed to hashed password",
		})
		return
	}

	hashed := models.User{Username: body.Username, Email: body.Email, Password: string(hash)}

	result := database.DB.Model(&body).Where("id = ?", id).Updates(&hashed)

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"message": "Email is already taken"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Edit user success", "data": hashed})
}

func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Param("id")

	if database.DB.Where("id = ?", id).Delete(&user).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete user success", "deletedUserId": id})
}

func Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	var user models.User

	if result := database.DB.Model(&user).Where("id = ?", userID).Update("token", ""); result.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to logout", "error": result.Error.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout success"})
}
