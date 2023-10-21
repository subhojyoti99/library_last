package controller

import (
	"Reference/database"
	"Reference/model"
	"Reference/util"
	"errors"
	"fmt"
	"net/http"

	// "os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// var validKey = os.Getenv("VALID_KEY")

func Register(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var input model.User

	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the role is valid
	if input.Role != "owner" && input.Role != "admin" && input.Role != "reader" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role specified"})
		return
	}

	// Set the role ID based on the input
	var roleID uint
	switch input.Role {
	case "owner":
		roleID = 1
	case "admin":
		roleID = 2
	case "reader":
		roleID = 3
	}

	owner := model.User{
		Name:          input.Name,
		ContactNumber: input.ContactNumber,
		Email:         input.Email,
		Password:      input.Password,
		RoleID:        roleID,
		Role:          input.Role,
		ValidKey:      input.ValidKey,
	}

	admin := model.User{
		Name:          input.Name,
		ContactNumber: input.ContactNumber,
		Email:         input.Email,
		Password:      input.Password,
		RoleID:        roleID,
		Role:          input.Role,
		LibID:         input.LibID,
		ValidKey:      input.ValidKey,
	}

	reader := model.User{
		Name:          input.Name,
		ContactNumber: input.ContactNumber,
		Email:         input.Email,
		Password:      input.Password,
		RoleID:        roleID,
		Role:          input.Role,
	}

	// if input.Role == "owner" && owner_reader.ValidKey == validKey {
	// 	fmt.Println("0000000", owner_reader.ValidKey, validKey)
	// 	context.JSON(http.StatusConflict, gin.H{"error": "Invalid Valid Key"})
	// 	return
	// }

	existingUser, err := model.GetUserByEmail(input.Email)
	_ = existingUser
	fmt.Println("hfiowehfiouhweiofhweiohf", existingUser)
	if err != nil {
		context.JSON(http.StatusConflict, gin.H{"error": "Duplicate email"})
		return
	}

	// existingLibrary, err := model.GetLibrary(int(admin.LibID))

	if input.Email != existingUser.Email {

		if input.Role == "owner" {
			database.DB.Create(&owner)
			context.JSON(http.StatusCreated, gin.H{"owner": owner})
		} else if input.Role == "admin" {
			database.DB.Create(&admin)
			context.JSON(http.StatusCreated, gin.H{"user": admin})
		} else {
			database.DB.Create(&reader)
			context.JSON(http.StatusCreated, gin.H{"reader": reader})
		}
	} else {
		context.JSON(http.StatusConflict, gin.H{"error": "User already registered."})
	}
	// database.DB.Save(&user)
	// context.JSON(http.StatusCreated, gin.H{"user": user})
	tx.Commit()
}

// User Login
func Login(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var input model.Login

	if err := context.ShouldBindJSON(&input); err != nil {
		var errorMessage string
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			validationError := validationErrors[0]
			if validationError.Tag() == "required" {
				errorMessage = fmt.Sprintf("%s not provided", validationError.Field())
			}
		}
		context.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
		return
	}

	user, err := model.GetUserByEmail(input.Email)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = user.ValidateUserPassword(input.Password)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jwt, err := util.GenerateJWT(user)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"token": jwt, "email": input.Email, "message": "Successfully logged in"})

	tx.Commit()
}

// get all users
func GetUsers(context *gin.Context) {
	var user []model.User
	err := model.GetUsers(&user)
	if err != nil {
		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, user)
}

// get user by id
func GetUser(context *gin.Context) {
	id, _ := strconv.Atoi(context.Param("id"))
	var user model.User
	err := model.GetUser(&user, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			context.AbortWithStatus(http.StatusNotFound)
			return
		}

		context.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, user)
}

// update user
func UpdateUser(c *gin.Context) {
	//var input model.Update
	var User model.User
	id, _ := strconv.Atoi(c.Param("id"))

	err := model.GetUser(&User, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.BindJSON(&User)
	err = model.UpdateUser(&User)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, User)
}
