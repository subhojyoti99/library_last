package controller

import (
	"Reference/database"
	"Reference/model"
	"Reference/util"
	"fmt"
	"net/http"
	"strconv"

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

	var approveIssue model.IssueRegistry

	err := context.ShouldBindJSON(&approveIssue)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requestID, err := strconv.Atoi(context.Param("request_id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	requests, err := model.FindRequestByID(requestID)

	bookDetails, err := model.FindBookByISBN(requests.BookISBN)
	if user.Role != "admin" && bookDetails.LibID != user.LibID {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. You can only update your own books."})
		return
	}

	issueApprove := model.IssueRegistry{
		ISBN:               requests.BookISBN,
		ReaderID:           requests.ReaderID,
		RequestID:          requests.ReqID,
		IssueApproverID:    user.ID,
		IssueStatus:        "approve",
		ExpectedReturnDate: approveIssue.ExpectedReturnDate,
	}

	bookDetails.AvailableCopies = bookDetails.AvailableCopies - 1
	database.DB.Create(&issueApprove)
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

	bookDetails, err := model.FindBookByISBN(returnIssue.ISBN)
	if user.Role != "admin" && bookDetails.LibID != user.LibID {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. You can only issue your own books."})
		return
	}

	issueGet.ReturnApproverID = user.ID
	issueGet.ReturnDate = returnIssue.ReturnDate
	issueGet.IssueStatus = "return"

	bookDetails.AvailableCopies = bookDetails.AvailableCopies + 1
	database.DB.Save(&issueGet)
	context.JSON(http.StatusOK, issueGet)

	database.DB.Save(&bookDetails)

	tx.Commit()
}
