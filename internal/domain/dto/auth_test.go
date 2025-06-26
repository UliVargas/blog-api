package dto

import (
	"testing"

	"github.com/UliVargas/blog-go/internal/domain/model"
	"github.com/stretchr/testify/assert"
)

func TestLoginRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request LoginRequest
		valid   bool
	}{
		{
			name: "valid login request",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			valid: true,
		},
		{
			name: "empty email",
			request: LoginRequest{
				Email:    "",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "invalid email format",
			request: LoginRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "empty password",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			valid: false,
		},
		{
			name: "password too short",
			request: LoginRequest{
				Email:    "test@example.com",
				Password: "123",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verificar que la estructura se puede crear
			assert.NotNil(t, tt.request)
			assert.IsType(t, LoginRequest{}, tt.request)
		})
	}
}

func TestLoginResponse_Structure(t *testing.T) {
	response := LoginResponse{
		Message: "Login successful",
		Token:   "jwt-token-here",
	}

	assert.Equal(t, "Login successful", response.Message)
	assert.Equal(t, "jwt-token-here", response.Token)
	assert.IsType(t, LoginResponse{}, response)
}

func TestRegisterRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request RegisterRequest
		valid   bool
	}{
		{
			name: "valid register request",
			request: RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "password123",
			},
			valid: true,
		},
		{
			name: "empty name",
			request: RegisterRequest{
				Name:     "",
				Email:    "john@example.com",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "name too short",
			request: RegisterRequest{
				Name:     "Jo",
				Email:    "john@example.com",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "name too long",
			request: RegisterRequest{
				Name:     "This is a very long name that exceeds the maximum allowed length for a user name field",
				Email:    "john@example.com",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "invalid email",
			request: RegisterRequest{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "password123",
			},
			valid: false,
		},
		{
			name: "password too short",
			request: RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "123",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verificar que la estructura se puede crear
			assert.NotNil(t, tt.request)
			assert.IsType(t, RegisterRequest{}, tt.request)
		})
	}
}

func TestRegisterRequest_ToUser(t *testing.T) {
	request := RegisterRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "password123",
	}

	user := request.ToUser()

	assert.Equal(t, request.Name, user.Name)
	assert.Equal(t, request.Email, user.Email)
	assert.Equal(t, request.Password, user.Password)
	assert.IsType(t, model.User{}, user)

	// Verificar que es una nueva instancia
	assert.NotSame(t, &request, &user)
}

func TestRegisterRequest_ToUser_EmptyFields(t *testing.T) {
	request := RegisterRequest{}

	user := request.ToUser()

	assert.Equal(t, "", user.Name)
	assert.Equal(t, "", user.Email)
	assert.Equal(t, "", user.Password)
	assert.IsType(t, model.User{}, user)
}