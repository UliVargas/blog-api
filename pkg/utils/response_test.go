package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSendSuccess(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		data     interface{}
		expected SuccessResponse
	}{
		{
			name:    "success with data",
			message: "Operation successful",
			data:    map[string]string{"id": "123"},
			expected: SuccessResponse{
				Message: "Operation successful",
				Data:    map[string]interface{}{"id": "123"},
			},
		},
		{
			name:    "success without data",
			message: "Operation completed",
			data:    nil,
			expected: SuccessResponse{
				Message: "Operation completed",
				Data:    nil,
			},
		},
		{
			name:    "success with string data",
			message: "User created",
			data:    "user123",
			expected: SuccessResponse{
				Message: "User created",
				Data:    "user123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.GET("/test", func(c *gin.Context) {
				SendSuccess(c, tt.message, tt.data)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			var response SuccessResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Message, response.Message)
			assert.Equal(t, tt.expected.Data, response.Data)
		})
	}
}

func TestSendCreated(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		data     interface{}
		expected SuccessResponse
	}{
		{
			name:    "created with data",
			message: "User created successfully",
			data:    map[string]string{"id": "456"},
			expected: SuccessResponse{
				Message: "User created successfully",
				Data:    map[string]interface{}{"id": "456"},
			},
		},
		{
			name:    "created without data",
			message: "Resource created",
			data:    nil,
			expected: SuccessResponse{
				Message: "Resource created",
				Data:    nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.POST("/test", func(c *gin.Context) {
				SendCreated(c, tt.message, tt.data)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusCreated, w.Code)

			var response SuccessResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Message, response.Message)
			assert.Equal(t, tt.expected.Data, response.Data)
		})
	}
}

func TestSendNoContent(t *testing.T) {
	router := setupTestRouter()
	router.DELETE("/test", func(c *gin.Context) {
		SendNoContent(c)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}
