// controller/controller_test.go
package controller

import (
	"Reference/database"
	"Reference/model"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/user/register", Register)

	database.ConnectDb()
	database.DB.AutoMigrate(&model.User{})
	defer database.DB.Migrator().DropTable(&model.User{})

	requestBody := `{
		"name": "John Doe",
		"email": "john.doe@example.com",
		"contact_number": "1234567890",
		"password": "password123",
		"role": "owner",
		"lib_id": 1
	}`

	req := performRequest(router, "POST", "/auth/user/register", requestBody)

	assert.Equal(t, http.StatusCreated, req.Code)
	assert.Contains(t, req.Body.String(), `"owner":`)
}

// performRequest is a helper function to perform HTTP requests for testing
func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}
