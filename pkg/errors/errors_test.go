package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		expected string
	}{
		{
			name: "with custom message",
			appError: &AppError{
				Err:        errors.New("base error"),
				Message:    "Custom error message",
				StatusCode: 400,
			},
			expected: "Custom error message",
		},
		{
			name: "without custom message",
			appError: &AppError{
				Err:        errors.New("base error message"),
				Message:    "",
				StatusCode: 500,
			},
			expected: "base error message",
		},
		{
			name: "empty message falls back to err",
			appError: &AppError{
				Err:        errors.New("fallback error"),
				StatusCode: 404,
			},
			expected: "fallback error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appError.Error()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAppError_Unwrap(t *testing.T) {
	baseErr := errors.New("base error")
	appErr := &AppError{
		Err:        baseErr,
		Message:    "Custom message",
		StatusCode: 400,
	}

	unwrapped := appErr.Unwrap()
	assert.Equal(t, baseErr, unwrapped)
}

func TestNewBadRequestError(t *testing.T) {
	baseErr := errors.New("validation failed")
	message := "Invalid input data"

	appErr := NewBadRequestError(baseErr, message)

	assert.Equal(t, baseErr, appErr.Err)
	assert.Equal(t, message, appErr.Message)
	assert.Equal(t, 400, appErr.StatusCode)
	assert.Equal(t, message, appErr.Error())
}

func TestNewNotFoundError(t *testing.T) {
	baseErr := errors.New("resource not found")
	message := "User not found"

	appErr := NewNotFoundError(baseErr, message)

	assert.Equal(t, baseErr, appErr.Err)
	assert.Equal(t, message, appErr.Message)
	assert.Equal(t, 404, appErr.StatusCode)
	assert.Equal(t, message, appErr.Error())
}

func TestNewConflictError(t *testing.T) {
	baseErr := errors.New("duplicate entry")
	message := "Email already exists"

	appErr := NewConflictError(baseErr, message)

	assert.Equal(t, baseErr, appErr.Err)
	assert.Equal(t, message, appErr.Message)
	assert.Equal(t, 409, appErr.StatusCode)
	assert.Equal(t, message, appErr.Error())
}

func TestNewInternalServerError(t *testing.T) {
	baseErr := errors.New("database connection failed")
	message := "Internal server error"

	appErr := NewInternalServerError(baseErr, message)

	assert.Equal(t, baseErr, appErr.Err)
	assert.Equal(t, message, appErr.Message)
	assert.Equal(t, 500, appErr.StatusCode)
	assert.Equal(t, message, appErr.Error())
}

func TestNewUnauthorizedError(t *testing.T) {
	baseErr := errors.New("invalid token")
	message := "Unauthorized access"

	appErr := NewUnauthorizedError(baseErr, message)

	assert.Equal(t, baseErr, appErr.Err)
	assert.Equal(t, message, appErr.Message)
	assert.Equal(t, 401, appErr.StatusCode)
	assert.Equal(t, message, appErr.Error())
}

func TestWrapDatabaseError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected error
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: nil,
		},
		{
			name:     "record not found",
			err:      errors.New("record not found"),
			expected: ErrUserNotFound,
		},
		{
			name:     "duplicate key with email",
			err:      errors.New("duplicate key constraint failed: email"),
			expected: ErrEmailExists,
		},
		{
			name:     "duplicate key with username",
			err:      errors.New("duplicate key constraint failed: username"),
			expected: ErrUsernameExists,
		},
		{
			name:     "UNIQUE constraint failed with email",
			err:      errors.New("UNIQUE constraint failed: users.email"),
			expected: ErrEmailExists,
		},
		{
			name:     "UNIQUE constraint failed with username",
			err:      errors.New("UNIQUE constraint failed: users.username"),
			expected: ErrUsernameExists,
		},
		{
			name:     "generic duplicate key",
			err:      errors.New("duplicate key constraint failed"),
			expected: ErrUserExists,
		},
		{
			name:     "foreign key constraint",
			err:      errors.New("foreign key constraint failed"),
			expected: ErrForeignKeyViolation,
		},
		{
			name:     "FOREIGN KEY constraint failed",
			err:      errors.New("FOREIGN KEY constraint failed"),
			expected: ErrForeignKeyViolation,
		},
		{
			name:     "connection refused",
			err:      errors.New("connection refused"),
			expected: ErrDatabaseConnection,
		},
		{
			name:     "no connection",
			err:      errors.New("no connection available"),
			expected: ErrDatabaseConnection,
		},
		{
			name: "generic database error",
			err:  errors.New("some database error"),
			// Para errores genéricos, debería envolver con ErrDatabaseOperation
			expected: nil, // Verificaremos que contiene ErrDatabaseOperation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapDatabaseError(tt.err)

			if tt.expected == nil && tt.err != nil {
				// Para errores genéricos, verificar que se envuelve con ErrDatabaseOperation
				assert.True(t, errors.Is(result, ErrDatabaseOperation))
				assert.Contains(t, result.Error(), "some database error")
			} else {
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestDomainErrors(t *testing.T) {
	// Test que los errores de dominio tienen los mensajes correctos
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{"ErrUserNotFound", ErrUserNotFound, "usuario no encontrado"},
		{"ErrUserExists", ErrUserExists, "el usuario ya existe"},
		{"ErrEmailExists", ErrEmailExists, "el email ya está registrado"},
		{"ErrUsernameExists", ErrUsernameExists, "el nombre de usuario ya está en uso"},
		{"ErrInvalidCredentials", ErrInvalidCredentials, "credenciales inválidas"},
		{"ErrUnauthorized", ErrUnauthorized, "no autorizado"},
		{"ErrInvalidInput", ErrInvalidInput, "datos de entrada inválidos"},
		{"ErrInvalidID", ErrInvalidID, "ID inválido"},
		{"ErrDatabaseConnection", ErrDatabaseConnection, "error de conexión con la base de datos"},
		{"ErrDatabaseOperation", ErrDatabaseOperation, "error en operación de base de datos"},
		{"ErrForeignKeyViolation", ErrForeignKeyViolation, "no se puede completar la operación debido a dependencias"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}