package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValidator(t *testing.T) {
	// Test que el validador se inicializa correctamente
	validator1 := GetValidator()
	assert.NotNil(t, validator1)

	// Test que siempre devuelve la misma instancia (singleton)
	validator2 := GetValidator()
	assert.Equal(t, validator1, validator2)
}

func TestFormatValidationErrors(t *testing.T) {
	// Test con validador real
	validator := GetValidator()

	// Estructura de prueba para validación básica
	type TestStruct struct {
		Name     string `validate:"required"`
		Email    string `validate:"email"`
		Password string `validate:"min=6"`
		Age      string `validate:"numeric"`
		Website  string `validate:"url"`
	}

	// Estructura para validaciones adicionales
	type ExtendedTestStruct struct {
		Code     string `validate:"len=4"`
		Name     string `validate:"alpha"`
		Username string `validate:"alphanum"`
		LongText string `validate:"max=10"`
	}

	tests := []struct {
		name     string
		data     interface{}
		expected map[string]string
	}{
		{
			name: "required field error",
			data: TestStruct{Name: "", Email: "test@example.com"},
			expected: map[string]string{
				"name": "Este campo es obligatorio",
			},
		},
		{
			name: "email validation error",
			data: TestStruct{Name: "Test", Email: "invalid-email"},
			expected: map[string]string{
				"email": "Debe ser un email válido",
			},
		},
		{
			name: "min length error",
			data: TestStruct{Name: "Test", Email: "test@example.com", Password: "123"},
			expected: map[string]string{
				"password": "Debe tener al menos 6 caracteres",
			},
		},
		{
			name: "numeric validation error",
			data: TestStruct{Name: "Test", Email: "test@example.com", Age: "not-a-number"},
			expected: map[string]string{
				"age": "Debe ser un número",
			},
		},
		{
			name: "url validation error",
			data: TestStruct{Name: "Test", Email: "test@example.com", Website: "not-a-url"},
			expected: map[string]string{
				"website": "Debe ser una URL válida",
			},
		},
		{
			name: "len validation error",
			data: ExtendedTestStruct{Code: "12345"}, // debe ser exactamente 4 caracteres
			expected: map[string]string{
				"code": "Debe tener exactamente 4 caracteres",
			},
		},
		{
			name: "alpha validation error",
			data: ExtendedTestStruct{Name: "Test123"}, // solo letras
			expected: map[string]string{
				"name": "Solo se permiten letras",
			},
		},
		{
			name: "alphanum validation error",
			data: ExtendedTestStruct{Username: "test@123"}, // solo letras y números
			expected: map[string]string{
				"username": "Solo se permiten letras y números",
			},
		},
		{
			name: "max length error",
			data: ExtendedTestStruct{LongText: "This is a very long text that exceeds the maximum"}, // más de 10 caracteres
			expected: map[string]string{
				"longtext": "No puede tener más de 10 caracteres",
			},
		},
	}

	// Test adicional para el caso default con una etiqueta personalizada
	type DefaultTestStruct struct {
		CustomField string `validate:"gte=5"` // etiqueta que no está en los casos específicos
	}

	tests = append(tests, struct {
		name     string
		data     interface{}
		expected map[string]string
	}{
		name: "default validation error",
		data: DefaultTestStruct{CustomField: "123"}, // menor que 5
		expected: map[string]string{
			"customfield": "Valor inválido",
		},
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Struct(tt.data)
			if err != nil {
				result := FormatValidationErrors(err)
				// Verificar que al menos contiene los errores esperados
				for expectedField, expectedMessage := range tt.expected {
					actualMessage, exists := result[expectedField]
					assert.True(t, exists, "Expected field %s to have validation error", expectedField)
					assert.Equal(t, expectedMessage, actualMessage)
				}
			}
		})
	}

	// Test con error no relacionado con validación
	t.Run("non-validation error", func(t *testing.T) {
		err := errors.New("some other error")
		result := FormatValidationErrors(err)
		assert.Empty(t, result)
	})
}

func TestCreateValidationErrorResponse(t *testing.T) {
	validator := GetValidator()

	// Estructura de prueba
	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"email"`
	}

	// Datos inválidos para generar errores
	invalidData := TestStruct{
		Name:  "",        // required error
		Email: "invalid", // email error
	}

	err := validator.Struct(invalidData)
	assert.Error(t, err)

	response := CreateValidationErrorResponse(err)

	assert.Equal(t, "Errores de validación", response.Error)
	assert.Contains(t, response.Errors, "name")
	assert.Contains(t, response.Errors, "email")
	assert.Equal(t, "Este campo es obligatorio", response.Errors["name"])
	assert.Equal(t, "Debe ser un email válido", response.Errors["email"])
}
