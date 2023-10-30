package controller

import (
	"Reference/database"
	"Reference/model"
	"Reference/util"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateRequest(context *gin.Context) {
	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)

	if user.Role != "reader" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. Only readers allowed."})
		return
	}

	var request model.RequestEvents

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookISBN := context.Param("book_isbn")

	book, err := model.FindBookByISBN(bookISBN)
	_ = book
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Book not found."})
		return
	}

	// Check if the role is valid
	if request.RequestType != "borrow" && request.RequestType != "return" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request specified"})
		return
	}

	// Check if the user has already made a request for this book
	existingRequest, err := model.FindRequestByBookAndReader(bookISBN, user.ID)
	if err == nil {
		if existingRequest.RequestType == "borrow" {
			fmt.Println("dtytdtdt", existingRequest.RequestType)
			existingRequest.RequestType = "return"
			existingRequest.RequestTypeID = 2
			database.DB.Save(&existingRequest)
			context.JSON(http.StatusOK, existingRequest)
			return
		}
		existingRequest.RequestType = "borrow"
		existingRequest.RequestTypeID = 1
		database.DB.Save(&existingRequest)
		context.JSON(http.StatusOK, existingRequest)
		return
	}

	// Set the role ID based on the input
	var RequestTypeID uint
	switch request.RequestType {
	case "borrow":
		RequestTypeID = 1
	case "return":
		RequestTypeID = 2
	}

	newRequest := model.RequestEvents{
		BookISBN:      bookISBN,
		ReaderID:      user.ID,
		ReadersEmail:  user.Email,
		RequestType:   request.RequestType,
		RequestTypeID: RequestTypeID,
		ApproverID:    book.AdminId,
	}
	if book.AvailableCopies == 0 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "No books available"})
		return
	}
	database.DB.Create(&newRequest)
	context.JSON(http.StatusCreated, newRequest)

	tx.Commit()
}

func GetRequests(context *gin.Context) {
	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)

	if user.Role != "admin" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. Only Admins allowed."})
		return
	}

	var requests []model.RequestEvents

	err := database.DB.Where("approver_id = ?", user.ID).Find(&requests).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"all_requests": requests})

	tx.Commit()
}

func GetRequestsAsReader(context *gin.Context) {
	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)

	if user.Role != "reader" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. Only Admins allowed."})
		return
	}

	var requests []model.RequestEvents

	err := database.DB.Where("reader_id = ?", user.ID).Find(&requests).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"My_requests": requests})

	tx.Commit()
}
