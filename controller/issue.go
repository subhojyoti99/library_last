package controller

import (
	"Reference/database"
	"Reference/model"
	"Reference/util"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ApproveIssue(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)

	var input model.IssueRegistry

	err := context.ShouldBindJSON(&input)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestID, err := strconv.Atoi(context.Param("request_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requests, _ := model.FindRequestByID(requestID)

	bookDetails, _ := model.FindBookByISBN(requests.BookISBN)
	if user.Role != "admin" && user.ID != bookDetails.AdminId {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. You can only update your own books."})
		return
	}

	var expReturnDate = time.Now().AddDate(0, 0, 20)

	issueApprove := model.IssueRegistry{
		ISBN:               requests.BookISBN,
		ReaderID:           requests.ReaderID,
		RequestID:          requests.ReqID,
		IssueApproverID:    user.ID,
		IssueStatus:        "approve",
		ExpectedReturnDate: expReturnDate,
	}

	database.DB.Create(&issueApprove)
	fmt.Println("--", issueApprove)
	bookDetails.AvailableCopies = bookDetails.AvailableCopies - 1
	context.JSON(http.StatusOK, issueApprove)

	database.DB.Save(&bookDetails)

	tx.Commit()
}

func ReturnIssue(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)

	var returnIssue model.IssueRegistry

	if err := context.ShouldBindJSON(&returnIssue); err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("7987896789", user)

	issue_id, _ := strconv.Atoi(context.Param("issue_id"))
	fmt.Println("7987896789", issue_id)

	issueGet, err := model.FindIssueByIssueID(issue_id)
	fmt.Println("7987896789", returnIssue)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Issue not found"})
		return
	}

	bookDetails, _ := model.FindBookByISBN(returnIssue.ISBN)
	if user.Role != "admin" && bookDetails.LibID != user.LibID {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. You can only issue your own books."})
		return
	}

	issueGet.ReturnApproverID = user.ID
	// issueGet.ReturnDate = returnIssue.ReturnDate
	issueGet.IssueStatus = "return"

	bookDetails.AvailableCopies = bookDetails.AvailableCopies + 1
	database.DB.Save(&issueGet)
	context.JSON(http.StatusOK, issueGet)

	database.DB.Save(&bookDetails)

	tx.Commit()
}

func GetAllIssues(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var issues []model.IssueRegistry

	user := util.CurrentUser(context)

	if user.Role != "admin" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. Only Admins allowed."})
		return
	}

	err := database.DB.Where("issue_approver_id = ?", user.ID).Find(&issues).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, gin.H{"All requests": issues})

	tx.Commit()
}
