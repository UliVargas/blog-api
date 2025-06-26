package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	services "github.com/UliVargas/blog-go/internal/application/service"
	"github.com/UliVargas/blog-go/internal/domain/model"
	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// MockUserService mocks the UserService for handler testing
type MockUserService struct {
	GetAllFunc  func() ([]model.User, error)
	GetByIDFunc func(id uint) (model.User, error)
	UpdateFunc  func(user model.User) (model.User, error)
	DeleteFunc  func(id uint) error
}

func (m *MockUserService) GetAll() ([]model.User, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return []model.User{}, nil
}

func (m *MockUserService) GetByID(id uint) (model.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return model.User{}, nil
}

func (m *MockUserService) Update(user model.User) (model.User, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return model.User{}, nil
}

func (m *MockUserService) Delete(id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

// NewUserHandlerWithMock creates a UserHandler with a mock service for testing
func NewUserHandlerWithMock() (*UserHandler, *MockUserService) {
	mockService := &MockUserService{}
	userHandler := &UserHandler{userService: mockService}
	return userHandler, mockService
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewUserHandler(t *testing.T) {
	mockService := &services.UserService{}
	userHandler := NewUserHandler(mockService)

	assert.NotNil(t, userHandler)
	assert.Equal(t, mockService, userHandler.userService)
}

func TestUserHandler_GetAll(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success - users found",
			mockSetup: func(m *MockUserService) {
				m.GetAllFunc = func() ([]model.User, error) {
					return []model.User{
						{ID: 1, Name: "John Doe", Email: "john@example.com"},
						{ID: 2, Name: "Jane Smith", Email: "jane@example.com"},
					}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"name":"John Doe","email":"john@example.com"},{"id":2,"name":"Jane Smith","email":"jane@example.com"}]`,
		},
		{
			name: "success - empty list",
			mockSetup: func(m *MockUserService) {
				m.GetAllFunc = func() ([]model.User, error) {
					return []model.User{}, nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[]`,
		},
		{
			name: "error - database error",
			mockSetup: func(m *MockUserService) {
				m.GetAllFunc = func() ([]model.User, error) {
					return nil, appErrors.NewInternalServerError(appErrors.ErrDatabaseOperation, "Database connection failed")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Database connection failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			userHandler, mockService := NewUserHandlerWithMock()
			tt.mockSetup(mockService)

			// Setup router and request
			router := setupRouter()
			router.GET("/users", userHandler.GetAll)

			// Create request
			req, _ := http.NewRequest("GET", "/users", nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse and compare JSON response
			if tt.expectedStatus == http.StatusOK {
				var actualUsers []model.User
				var expectedUsers []model.User
				err := json.Unmarshal(w.Body.Bytes(), &actualUsers)
				assert.NoError(t, err)
				err = json.Unmarshal([]byte(tt.expectedBody), &expectedUsers)
				assert.NoError(t, err)
				assert.Equal(t, expectedUsers, actualUsers)
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	tests := []struct {
		name           string
		urlParam       string
		mockSetup      func(*MockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:     "success - user found",
			urlParam: "1",
			mockSetup: func(m *MockUserService) {
				m.GetByIDFunc = func(id uint) (model.User, error) {
					if id == 1 {
						return model.User{ID: 1, Name: "John Doe", Email: "john@example.com"}, nil
					}
					return model.User{}, appErrors.NewNotFoundError(appErrors.ErrUserNotFound, "User not found")
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"name":"John Doe","email":"john@example.com"}`,
		},
		{
			name:           "error - invalid ID format",
			urlParam:       "invalid",
			mockSetup:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"ID inválido"}`,
		},
		{
			name:     "error - user not found",
			urlParam: "999",
			mockSetup: func(m *MockUserService) {
				m.GetByIDFunc = func(id uint) (model.User, error) {
					return model.User{}, appErrors.NewNotFoundError(appErrors.ErrUserNotFound, "User not found")
				}
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"User not found"}`,
		},
		{
			name:     "error - database error",
			urlParam: "1",
			mockSetup: func(m *MockUserService) {
				m.GetByIDFunc = func(id uint) (model.User, error) {
					return model.User{}, appErrors.NewInternalServerError(appErrors.ErrDatabaseOperation, "Database connection failed")
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Database connection failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			userHandler, mockService := NewUserHandlerWithMock()
			tt.mockSetup(mockService)

			// Setup router and request
			router := setupRouter()
			router.GET("/users/:id", userHandler.GetByID)

			// Create request
			req, _ := http.NewRequest("GET", "/users/"+tt.urlParam, nil)
			w := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var actualUser model.User
				var expectedUser model.User
				err := json.Unmarshal(w.Body.Bytes(), &actualUser)
				assert.NoError(t, err)
				err = json.Unmarshal([]byte(tt.expectedBody), &expectedUser)
				assert.NoError(t, err)
				assert.Equal(t, expectedUser, actualUser)
			} else {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}
