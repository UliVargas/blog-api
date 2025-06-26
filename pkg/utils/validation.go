package utils

import (
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	validatorInstance *validator.Validate
	once              sync.Once
)

// GetValidator devuelve una instancia del validador
func GetValidator() *validator.Validate {
	once.Do(func() {
		validatorInstance = validator.New()
	})
	return validatorInstance
}

// ValidationErrorResponse representa la estructura de respuesta para errores de validació
type ValidationErrorResponse struct {
	Error  string            `json:"error"`
	Errors map[string]string `json:"errors"`
}

// FormatValidationErrors convierte los errores de validación en mensajes más amigables
func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			fieldName := strings.ToLower(fieldError.Field())

			switch fieldError.Tag() {
			case "required":
				errors[fieldName] = "Este campo es obligatorio"
			case "min":
				errors[fieldName] = "Debe tener al menos " + fieldError.Param() + " caracteres"
			case "max":
				errors[fieldName] = "No puede tener más de " + fieldError.Param() + " caracteres"
			case "email":
				errors[fieldName] = "Debe ser un email válido"
			case "len":
				errors[fieldName] = "Debe tener exactamente " + fieldError.Param() + " caracteres"
			case "numeric":
				errors[fieldName] = "Debe ser un número"
			case "alpha":
				errors[fieldName] = "Solo se permiten letras"
			case "alphanum":
				errors[fieldName] = "Solo se permiten letras y números"
			case "url":
				errors[fieldName] = "Debe ser una URL válida"
			default:
				errors[fieldName] = "Valor inválido"
			}
		}
	}

	return errors
}

// CreateValidationErrorResponse crea una respuesta estándar para errores de validación
func CreateValidationErrorResponse(err error) ValidationErrorResponse {
	return ValidationErrorResponse{
		Error:  "Errores de validación",
		Errors: FormatValidationErrors(err),
	}
}
