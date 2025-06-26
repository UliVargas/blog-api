package errors

import (
	"errors"
	"fmt"
	"strings"
)

// Errores de dominio personalizados
var (
	// Errores de usuario
	ErrUserNotFound    = errors.New("usuario no encontrado")
	ErrUserExists      = errors.New("el usuario ya existe")
	ErrEmailExists     = errors.New("el email ya está registrado")
	ErrUsernameExists  = errors.New("el nombre de usuario ya está en uso")
	
	// Errores de autenticación
	ErrInvalidCredentials = errors.New("credenciales inválidas")
	ErrUnauthorized      = errors.New("no autorizado")
	
	// Errores de validación
	ErrInvalidInput = errors.New("datos de entrada inválidos")
	ErrInvalidID    = errors.New("ID inválido")
	
	// Errores de base de datos
	ErrDatabaseConnection = errors.New("error de conexión con la base de datos")
	ErrDatabaseOperation  = errors.New("error en operación de base de datos")
	ErrForeignKeyViolation = errors.New("no se puede completar la operación debido a dependencias")
)

// AppError representa un error de aplicación con código HTTP
type AppError struct {
	Err        error
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Constructores de errores con códigos HTTP
func NewBadRequestError(err error, message string) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: 400,
	}
}

func NewNotFoundError(err error, message string) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: 404,
	}
}

func NewConflictError(err error, message string) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: 409,
	}
}

func NewInternalServerError(err error, message string) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: 500,
	}
}

func NewUnauthorizedError(err error, message string) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: 401,
	}
}

// WrapDatabaseError convierte errores de GORM en errores de dominio
func WrapDatabaseError(err error) error {
	if err == nil {
		return nil
	}
	
	errorMsg := err.Error()
	
	// Detectar errores específicos de GORM
	if errorMsg == "record not found" {
		return ErrUserNotFound
	}
	
	// Detectar errores de duplicación
	if strings.Contains(errorMsg, "duplicate key") || strings.Contains(errorMsg, "UNIQUE constraint failed") {
		if strings.Contains(errorMsg, "email") {
			return ErrEmailExists
		}
		if strings.Contains(errorMsg, "username") {
			return ErrUsernameExists
		}
		return ErrUserExists
	}
	
	// Detectar errores de clave foránea
	if strings.Contains(errorMsg, "foreign key") || strings.Contains(errorMsg, "FOREIGN KEY constraint failed") {
		return ErrForeignKeyViolation
	}
	
	// Detectar errores de conexión
	if strings.Contains(errorMsg, "connection refused") || strings.Contains(errorMsg, "no connection") {
		return ErrDatabaseConnection
	}
	
	// Para otros errores, envolver como error de operación de base de datos
	return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
}