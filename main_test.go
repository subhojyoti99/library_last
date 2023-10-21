// main_test.go
package main

import (
	"Reference/controller"
	"Reference/database"
	"Reference/model"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Load environment variables
	err := godotenv.Load(".env.test")
	if err != nil {
		panic("Error loading .env.test file")
	}

	// Connect to test database
	database.ConnectDb()
	database.DB.AutoMigrate(&model.User{})
	database.DB.AutoMigrate(&model.Library{})
	database.DB.AutoMigrate(&model.BookInventory{})
	database.DB.AutoMigrate(&model.RequestEvents{})
	database.DB.AutoMigrate(&model.IssueRegistry{})

	// Run the tests
	exitCode := m.Run()

	// Clean up after tests
	database.DB.Migrator().DropTable(&model.User{}, &model.Library{}, &model.BookInventory{}, &model.RequestEvents{}, &model.IssueRegistry{})
	os.Exit(exitCode)
}

func TestRegister(t *testing.T) {
	// Setup a test Gin context
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/user/register", controller.Register)

	// Define the JSON request body
	requestBody := `{
		"name": "John Doe",
		"email": "john.doe@example.com",
		"contact_number": "1234567890",
		"password": "password123",
		"role": "owner",
		"lib_id": 1
	}`

	// Perform a POST request with the test context
	req := performRequest(router, "POST", "/auth/user/register", requestBody)

	// Assertions
	assert.Equal(t, 201, req.Code)
	assert.Contains(t, req.Body.String(), `"owner":`)
}

func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
