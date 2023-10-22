package model

import (
	"Reference/database"
	"fmt"
	"time"
)

type IssueRegistry struct {
	IssueID            uint      `gorm:"primary_key;auto_increment" json:"issue_id"`
	ISBN               string    `json:"isbn"`
	RequestID          uint      `json:"request_id"`
	ReaderID           uint      `json:"reader_id"`
	IssueApproverID    uint      `json:"issue_approver_id"`
	IssueStatus        string    `json:"issue_status"`
	CreatedAt          time.Time `json:"issue_date"`
	ExpectedReturnDate time.Time `json:"expected_return_date"`
	UpdatedAt          time.Time `json:"return_date"`
	ReturnApproverID   uint      `json:"return_approver_id"`
}

func (issue *IssueRegistry) Save() (*IssueRegistry, error) {
	err := database.DB.Create(&issue).Error
	if err != nil {
		return &IssueRegistry{}, err
	}
	return issue, err
}

func FindIssueByIssueID(issue_id int) (IssueRegistry, error) {
	var issue IssueRegistry
	err := database.DB.Where("issue_id = ?", issue_id).First(&issue).Error
	fmt.Println("12345678", err, issue)
	if err != nil {
		return IssueRegistry{}, err
	}
	return issue, nil
}
