package model

import (
	"Reference/database"
)

// Book represents a book entity.
type BookInventory struct {
	ISBN            string  `gorm:"primary_key" json:"isbn"`
	AdminId         uint    `json:"admin_id"`
	LibID           uint    `json:"library_id"`
	Title           string  `json:"title"`
	Authors         string  `json:"authors"`
	Publisher       string  `json:"publisher"`
	Version         string  `json:"version"`
	TotalCopies     int     `json:"total_copies"`
	AvailableCopies int     `json:"available_copies"`
	Library         Library `gorm:"foreignkey:LibID"`
}

func (book *BookInventory) Save() (*BookInventory, error) {
	err := database.DB.Create(&book).Error
	if err != nil {
		return &BookInventory{}, err
	}
	return book, err
}

func FindBookByISBN(isbn string) (BookInventory, error) {
	var book BookInventory
	err := database.DB.Where("isbn = ?", isbn).First(&book).Error
	if err != nil {
		return BookInventory{}, err
	}
	return book, nil
}

func SearchBooks(query string) ([]BookInventory, error) {
	var books []BookInventory

	// Assuming 'Title' and 'Authors' are the fields you want to search
	err := database.DB.Where("title LIKE ? OR authors LIKE ?", "%"+query+"%", "%"+query+"%").Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}
