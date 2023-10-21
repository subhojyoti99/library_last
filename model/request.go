package model

import (
	"Reference/database"
	"time"
)

type RequestEvents struct {
	ReqID         uint      `gorm:"primary_key;auto_increment" json:"request_id"`
	BookISBN      string    `json:"book_isbn"`
	ReaderID      uint      `json:"reader_id"`
	ReadersEmail  string    `json:"readers_email"`
	CreatedAt     time.Time `json:"request_date"`
	UpdatedAt     time.Time `json:"approval_date"`
	ApproverID    uint      `json:"approver_id"`
	RequestType   string    `json:"request_type"`
	RequestTypeID uint      `json:"request_type_id"`
}

func (request *RequestEvents) Save() (*RequestEvents, error) {
	err := database.DB.Create(&request).Error
	if err != nil {
		return &RequestEvents{}, err
	}
	return request, nil
}

func FindRequestByID(request_id int) (RequestEvents, error) {
	var request RequestEvents
	err := database.DB.Where("req_id = ?", request_id).First(&request).Error
	if err != nil {
		return RequestEvents{}, err
	}
	return request, nil
}

func FindRequestBookByISBN(book_isbn string) (RequestEvents, error) {
	var bookRequest RequestEvents
	err := database.DB.Where("book_isbn = ?", book_isbn).First(&bookRequest).Error
	if err != nil {
		return RequestEvents{}, err
	}
	return bookRequest, nil
}

func FindRequestByBookAndReader(bookISBN string, readerID uint) (RequestEvents, error) {
	var request RequestEvents
	err := database.DB.Where("book_isbn = ? AND reader_id = ?", bookISBN, readerID).First(&request).Error
	return request, err
}
