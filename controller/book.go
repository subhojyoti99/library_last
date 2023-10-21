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

func AddBookInventory(context *gin.Context) {

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

	var book model.BookInventory

	err := context.ShouldBindJSON(&book)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	library, err := model.GetLibrary(int(user.LibID))
	fmt.Println("======", library)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Library not found."})
		return
	}

	// Check if a book with the same ISBN already exists
	existingBook, err := model.FindBookByISBN(book.ISBN)
	_ = existingBook
	if err == nil {
		context.JSON(http.StatusConflict, gin.H{"error": "Duplicate ISBN"})
		return
	}

	newBook := model.BookInventory{
		ISBN:            book.ISBN,
		LibID:           uint(user.LibID),
		AdminId:         user.ID,
		Title:           book.Title,
		Authors:         book.Authors,
		Publisher:       book.Publisher,
		Version:         book.Version,
		TotalCopies:     book.TotalCopies,
		AvailableCopies: book.AvailableCopies,
		Library:         library,
	}

	addedBook, err := newBook.Save()

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"book": addedBook})

	tx.Commit()
}

func GetAllBook(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var books []model.BookInventory

	if err := database.DB.Find(&books).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books."})
		return
	}
	context.JSON(http.StatusOK, gin.H{"Books": books})

	tx.Commit()
}

func GetBookByLibraryID(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var books []model.BookInventory

	libraryID, _ := strconv.Atoi(context.Param("library_id"))

	library, err := model.GetLibrary(libraryID)
	_ = library
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Library not found."})
		return
	}

	err = database.DB.Where("lib_id = ?", libraryID).Find(&books).Error
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for i := 0; i < len(books); i++ {
		books[i].Library = library
	}

	context.JSON(http.StatusOK, books)

	tx.Commit()
}

// Update book
func UpdateBook(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)

	var input model.BookInventory
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bookISBN := context.Param("isbn")

	book, err := model.FindBookByISBN(bookISBN)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	if book.LibID != user.LibID {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. You can only update your own books."})
		return
	}

	book.AdminId = user.ID
	book.Title = input.Title
	book.Authors = input.Authors
	book.Publisher = input.Publisher
	book.Version = input.Version
	book.TotalCopies = input.TotalCopies
	book.AvailableCopies = input.AvailableCopies

	database.DB.Save(&book)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, book)

	tx.Commit()
}

func DeleteBook(context *gin.Context) {

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := util.CurrentUser(context)
	if user.Role != "admin" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. Only Owners allowed."})
		return
	}

	bookISBN := context.Param("isbn")

	book, err := model.FindBookByISBN(bookISBN)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Book not found"})
		return
	}

	if user.LibID == book.LibID {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "Access denied. Only Owners allowed."})
		return
	}

	database.DB.Delete(book)
	context.JSON(http.StatusOK, "Book has been Deleted")

	tx.Commit()
}

func SearchBooks(c *gin.Context) {
	query := c.Query("query") // Get the search query from the URL parameter
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}
	// Perform the search
	results, err := model.SearchBooks(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, results)
}
