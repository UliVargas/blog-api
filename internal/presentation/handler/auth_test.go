package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	services "github.com/UliVargas/blog-go/internal/application/service"
	"github.com/UliVargas/blog-go/internal/domain/dto"
	"github.com/UliVargas/blog-go/internal/domain/model"
	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"github.com/stretchr/testify/assert"
)

// MockAuthService mocks the AuthService for handler testing
type MockAuthService struct {
	LoginFunc    func(email, password string) (string, error)
	RegisterFunc func(user model.User) error
}

func (m *MockAuthService) Login(email, password string) (string, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(email, password)
	}
	return "mock-token", nil
}

func (m *MockAuthService) Register(user model.User) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(user)
	}
	return nil
}

// NewAuthHandlerWithMock creates an AuthHandler with a mock service for testing
func NewAuthHandlerWithMock() (*AuthHandler, *MockAuthService) {
	mockService := &MockAuthService{}
	authHandler := &AuthHandler{authService: mockService}
	return authHandler, mockService
}

func TestNewAuthHandler(t *testing.T) {
	mockService := &services.AuthService{}
	authHandler := NewAuthHandler(mockService)

	assert.NotNil(t, authHandler)
	assert.Equal(t, mockService, authHandler.authService)
}

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success - valid credentials",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.LoginFunc = func(email, password string) (string, error) {
					return "jwt-token-123", nil
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"Inicio de sesión exitoso","token":"jwt-token-123"}`,
		},
		{
			name: "error - invalid credentials",
			requestBody: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *MockAuthService) {
				m.LoginFunc = func(email, password string) (string, error) {
					return "", appErrors.ErrInvalidCredentials
				}
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Credenciales inválidas"}`,
		},
		{
			name:        "error - invalid JSON",
			requestBody: `{"email":"invalid-json"`,
			mockSetup: func(m *MockAuthService) {
				// No setup needed for this test
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Datos inválidos"}`,
		},
		{
			name: "error - missing email",
			requestBody: dto.LoginRequest{
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				// No setup needed for this test
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Datos de validación incorrectos","errors":{"email":"Este campo es obligatorio"}}`,
		},
		{
			name: "error - missing password",
			requestBody: dto.LoginRequest{
				Email: "test@example.com",
			},
			mockSetup: func(m *MockAuthService) {
				// No setup needed for this test
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Datos de validación incorrectos","errors":{"password":"Este campo es obligatorio"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authHandler, mockService := NewAuthHandlerWithMock()
			tt.mockSetup(mockService)

			router := setupRouter()
			router.POST("/login", authHandler.Login)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockSetup      func(*MockAuthService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success - valid user data",
			requestBody: dto.RegisterRequest{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.RegisterFunc = func(user model.User) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"Usuario creado exitosamente"}`,
		},
		{
			name: "error - email already exists",
			requestBody: dto.RegisterRequest{
				Name:     "Test User",
				Email:    "existing@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				m.RegisterFunc = func(user model.User) error {
					return appErrors.ErrEmailExists
				}
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"El email ya está registrado"}`,
		},
		{
			name:        "error - invalid JSON",
			requestBody: `{"name":"invalid-json"`,
			mockSetup: func(m *MockAuthService) {
				// No setup needed for this test
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Datos inválidos proporcionados"}`,
		},
		{
			name: "error - missing name",
			requestBody: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				// No setup needed for this test
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Datos de validación incorrectos","errors":{"name":"Este campo es obligatorio"}}`,
		},
		{
			name: "error - missing email",
			requestBody: dto.RegisterRequest{
				Name:     "Test User",
				Password: "password123",
			},
			mockSetup: func(m *MockAuthService) {
				// No setup needed for this test
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Datos de validación incorrectos","errors":{"email":"Este campo es obligatorio"}}`,
		},
		{
			name: "error - missing password",
			requestBody: dto.RegisterRequest{
				Name:  "Test User",
				Email: "test@example.com",
			},
			mockSetup: func(m *MockAuthService) {
				// No setup needed for this test
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Datos de validación incorrectos","errors":{"password":"Este campo es obligatorio"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authHandler, mockService := NewAuthHandlerWithMock()
			tt.mockSetup(mockService)

			router := setupRouter()
			router.POST("/register", authHandler.Register)

			var body []byte
			if str, ok := tt.requestBody.(string); ok {
				body = []byte(str)
			} else {
				body, _ = json.Marshal(tt.requestBody)
			}

			req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
