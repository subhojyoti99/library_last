// controller/library_controller.go

package controller

import (
	"Reference/database"
	"Reference/model"
	"Reference/util"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateLibrary(context *gin.Context) {
	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)

	if user.Role != "owner" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. Only Owners allowed."})
		return
	}

	var library model.Library

	err := context.ShouldBindJSON(&library)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingLibrary, err := model.FindLibraryByName(library.Name)
	_ = existingLibrary
	if err == nil {
		context.JSON(http.StatusConflict, gin.H{"error": "Duplicate Name"})
		return
	}

	newLibrary := model.Library{
		Name:   library.Name,
		UserId: user.ID, // Associate the library with the current user
	}

	database.DB.Create(&newLibrary)

	user.LibID = newLibrary.Id
	database.DB.Save(&user)
	database.DB.Save(&newLibrary)

	context.JSON(http.StatusCreated, gin.H{"library": newLibrary})

	tx.Commit()
}

func GetAllLibraries(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var libraries []model.Library

	if err := database.DB.Find(&libraries).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch libraries"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"libraries": libraries})

	tx.Commit()
}

// Update library
func UpdateLibrary(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)

	var input model.Library
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	libraryID, _ := strconv.Atoi(context.Param("id"))

	library, err := model.GetLibrary(libraryID)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Library not found"})
		return
	}

	if library.UserId != user.ID {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. You can only update your own libraries."})
		return
	}

	library.Name = input.Name

	err = model.UpdateLibrary(&library)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, library)

	tx.Commit()
}

func DeleteLibrary(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)
	if user.Role != "owner" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. Only Owners allowed."})
		return
	}

	libraryID, _ := strconv.Atoi(context.Param("id"))

	library, err := model.GetLibrary(libraryID)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Library not found"})
		return
	}

	if user.LibID == library.Id {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. Only Owners allowed."})
		return
	}

	database.DB.Delete(library)
	context.JSON(http.StatusOK, "Book has been Deleted")

	tx.Commit()
}
