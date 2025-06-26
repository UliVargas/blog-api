package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "nil error",
			err:            nil,
			expectedStatus: 0, // No response expected
		},
		{
			name:           "AppError with custom message",
			err:            appErrors.NewBadRequestError(errors.New("base error"), "Custom message"),
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Custom message",
		},
		{
			name:           "AppError without custom message",
			err:            appErrors.NewNotFoundError(errors.New("not found"), ""),
			expectedStatus: http.StatusNotFound,
			expectedError:  "not found",
		},
		{
			name:           "ErrUserNotFound",
			err:            appErrors.ErrUserNotFound,
			expectedStatus: http.StatusNotFound,
			expectedError:  "Usuario no encontrado",
		},
		{
			name:           "ErrEmailExists",
			err:            appErrors.ErrEmailExists,
			expectedStatus: http.StatusConflict,
			expectedError:  "El email ya está registrado",
		},
		{
			name:           "ErrUsernameExists",
			err:            appErrors.ErrUsernameExists,
			expectedStatus: http.StatusConflict,
			expectedError:  "El nombre de usuario ya está en uso",
		},
		{
			name:           "ErrUserExists",
			err:            appErrors.ErrUserExists,
			expectedStatus: http.StatusConflict,
			expectedError:  "El usuario ya existe",
		},
		{
			name:           "ErrInvalidCredentials",
			err:            appErrors.ErrInvalidCredentials,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Credenciales inválidas",
		},
		{
			name:           "ErrUnauthorized",
			err:            appErrors.ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "No autorizado",
		},
		{
			name:           "ErrInvalidInput",
			err:            appErrors.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Datos de entrada inválidos",
		},
		{
			name:           "ErrInvalidID",
			err:            appErrors.ErrInvalidID,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "ID inválido",
		},
		{
			name:           "ErrDatabaseConnection",
			err:            appErrors.ErrDatabaseConnection,
			expectedStatus: http.StatusServiceUnavailable,
			expectedError:  "Error de conexión con la base de datos",
		},
		{
			name:           "ErrForeignKeyViolation",
			err:            appErrors.ErrForeignKeyViolation,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "No se puede completar la operación debido a dependencias",
		},
		{
			name:           "ErrDatabaseOperation",
			err:            appErrors.ErrDatabaseOperation,
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Error en operación de base de datos",
		},
		{
			name:           "Generic error",
			err:            errors.New("some generic error"),
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Error interno del servidor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.GET("/test", func(c *gin.Context) {
				HandleError(c, tt.err)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			if tt.err == nil {
				// Para errores nil, no debería haber respuesta
				assert.Equal(t, http.StatusOK, w.Code) // Gin devuelve 200 por defecto si no se establece status
				return
			}

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedError, response.Error)

			// Para errores genéricos, verificar que se incluyen los detalles
			if tt.name == "Generic error" {
				assert.Equal(t, "some generic error", response.Details)
			}
		})
	}
}

func TestHandleValidationError(t *testing.T) {
	// Usar el validador real para generar errores de validación
	validator := GetValidator()

	// Estructura de prueba
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"email"`
	}

	// Datos inválidos para generar errores
	invalidData := TestStruct{
		Name:  "",           // required error
		Email: "invalid",    // email error
	}

	validationErr := validator.Struct(invalidData)
	assert.Error(t, validationErr)

	router := setupTestRouter()
	router.POST("/test", func(c *gin.Context) {
		HandleValidationError(c, validationErr)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ValidationErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Datos de validación incorrectos", response.Error)
	assert.Contains(t, response.Errors, "name")
	assert.Contains(t, response.Errors, "email")
	assert.Equal(t, "Este campo es obligatorio", response.Errors["name"])
	assert.Equal(t, "Debe ser un email válido", response.Errors["email"])
}

func TestHandleBadRequest(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{
			name:    "custom bad request message",
			message: "Invalid request format",
		},
		{
			name:    "empty message",
			message: "",
		},
		{
			name:    "long message",
			message: "This is a very long error message that describes what went wrong with the request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.POST("/test", func(c *gin.Context) {
				HandleBadRequest(c, tt.message)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response ErrorResponse
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.message, response.Error)
			assert.Empty(t, response.Details) // Details should be empty for bad request
		})
	}
}
