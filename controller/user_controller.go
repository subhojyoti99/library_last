package controller

import (
	"Reference/database"
	"Reference/model"
	"Reference/util"
	"errors"
	"fmt"
	"net/http"

	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func Register(context *gin.Context) {
	ownerKey := os.Getenv("OWNER_VALID_KEY")
	adminKey := os.Getenv("ADMIN_VALID_KEY")

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

	existingUser, _ := model.GetUserByEmail(input.Email)

	if input.Email == existingUser.Email {
		context.JSON(http.StatusConflict, gin.H{"error": "User already registered."})
		return
	}
	fmt.Println("hfiowehfiouhwe--iofhweiohf", ownerKey)
	if input.Role == "owner" && input.ValidKey == ownerKey {
		database.DB.Create(&owner)
		context.JSON(http.StatusCreated, gin.H{"owner": owner})
	} else if input.Role == "admin" && input.ValidKey == adminKey {
		database.DB.Create(&admin)
		context.JSON(http.StatusCreated, gin.H{"admin": admin})
	} else if input.Role == "reader" {
		database.DB.Create(&reader)
		context.JSON(http.StatusCreated, gin.H{"reader": reader})

	} else {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Enter correct valid key.."})
	}
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

	context.JSON(http.StatusOK, gin.H{"token": jwt, "user": user, "message": "Successfully logged in"})

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
