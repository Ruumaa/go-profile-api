package controllers

import (
	"go-web/database"
	"go-web/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllPhotos(c *gin.Context) {
	var photos []models.Photo

	database.DB.Find(&photos)

	c.JSON(http.StatusOK, gin.H{"message": "GET photos success", "data": photos})
}

func CreatePhoto(c *gin.Context) {
	var photo models.Photo

	if err := c.ShouldBindJSON(&photo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if result := database.DB.Create(&photo); result.RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Id not found", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Create photo success", "data": photo})
}

func EditPhoto(c *gin.Context) {
	var photo models.Photo
	id := c.Param("photoId")

	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Periksa apakah photo dengan ID yang diberikan ada
	if err := database.DB.First(&photo, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Photo not found"})
		return
	}

	// Periksa apakah pengguna yang mengedit sesuai dengan pembuat photo
	if photo.UserID != userID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "You are not allowed to edit this photo"})
		return
	}

	if err := c.ShouldBindJSON(&photo); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if result := database.DB.Model(&photo).Where("id = ?", id).Updates(&photo); result.Error != nil {
		if result.RowsAffected == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Id not found", "error": result.Error.Error()})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to edit Photo ", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Edit photo success", "data": photo})

}

func DeletePhoto(c *gin.Context) {
	var photo models.Photo
	id := c.Param("photoId")

	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	// Periksa apakah photo dengan ID yang diberikan ada
	if err := database.DB.First(&photo, id).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"message": "Photo not found"})
		return
	}

	// Periksa apakah pengguna yang mengedit sesuai dengan pembuat photo
	if photo.UserID != userID {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "You are not allowed to delete this photo"})
		return
	}

	if database.DB.Where("id = ?", id).Delete(&photo).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delete photo success", "deletedUserId": id})
}

func Validate(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userID": userID})
}
