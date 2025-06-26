package service

import (
	"errors"
	"os"
	"testing"

	"github.com/UliVargas/blog-go/internal/domain/model"
	"github.com/UliVargas/blog-go/internal/domain/repository"
	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepositoryAuth mocks the UserRepository for auth testing
type MockUserRepositoryAuth struct {
	mock.Mock
}

func (m *MockUserRepositoryAuth) GetAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepositoryAuth) GetByID(id uint) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepositoryAuth) GetByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepositoryAuth) Create(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepositoryAuth) Update(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepositoryAuth) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// NewAuthServiceWithMock creates an AuthService with a mock repository for testing
func NewAuthServiceWithMock() (*AuthService, *MockUserRepositoryAuth) {
	mockRepo := &MockUserRepositoryAuth{}
	service := NewAuthService(mockRepo)
	return service, mockRepo
}

func TestNewAuthService(t *testing.T) {
	tests := []struct {
		name     string
		userRepo repository.UserRepositoryInterface
		wantNil  bool
	}{
		{
			name:     "success - valid repository",
			userRepo: &MockUserRepositoryAuth{},
			wantNil:  false,
		},
		{
			name:     "success - nil repository",
			userRepo: nil,
			wantNil:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewAuthService(tt.userRepo)
			if tt.wantNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
				assert.Equal(t, tt.userRepo, result.userRepo)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	// Create a hashed password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	tests := []struct {
		name         string
		email        string
		password     string
		mockSetup    func(*MockUserRepositoryAuth)
		wantToken    bool
		wantError    error
		errorMessage string
	}{
		{
			name:     "success - valid credentials",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(m *MockUserRepositoryAuth) {
				m.On("GetByEmail", "test@example.com").Return(
					model.User{
						ID:       1,
						Email:    "test@example.com",
						Password: string(hashedPassword),
					},
					nil,
				)
			},
			wantToken: true,
			wantError: nil,
		},
		{
			name:     "error - user not found",
			email:    "notfound@example.com",
			password: "password123",
			mockSetup: func(m *MockUserRepositoryAuth) {
				m.On("GetByEmail", "notfound@example.com").Return(
					model.User{},
					appErrors.ErrUserNotFound,
				)
			},
			wantToken: false,
			wantError: appErrors.ErrInvalidCredentials,
		},
		{
			name:     "error - database error",
			email:    "test@example.com",
			password: "password123",
			mockSetup: func(m *MockUserRepositoryAuth) {
				m.On("GetByEmail", "test@example.com").Return(
					model.User{},
					errors.New("database connection error"),
				)
			},
			wantToken: false,
			wantError: errors.New("database connection error"),
		},
		{
			name:     "error - invalid password",
			email:    "test@example.com",
			password: "wrongpassword",
			mockSetup: func(m *MockUserRepositoryAuth) {
				m.On("GetByEmail", "test@example.com").Return(
					model.User{
						ID:       1,
						Email:    "test@example.com",
						Password: string(hashedPassword),
					},
					nil,
				)
			},
			wantToken: false,
			wantError: appErrors.ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			service, mockRepo := NewAuthServiceWithMock()
			tt.mockSetup(mockRepo)

			// Execute
			token, err := service.Login(tt.email, tt.password)

			// Assert
			if tt.wantError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantError.Error(), err.Error())
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				if tt.wantToken {
					assert.NotEmpty(t, token)
					// Verify it's a valid JWT token structure
					_, _, err := new(jwt.Parser).ParseUnverified(token, jwt.MapClaims{})
					assert.NoError(t, err, "Token should be a valid JWT")
				} else {
					assert.Empty(t, token)
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name      string
		user      model.User
		mockSetup func(*MockUserRepositoryAuth)
		wantError error
	}{
		{
			name: "success - new user registration",
			user: model.User{
				Email:    "newuser@example.com",
				Password: "password123",
				Name:     "New User",
			},
			mockSetup: func(m *MockUserRepositoryAuth) {
				// User doesn't exist
				m.On("GetByEmail", "newuser@example.com").Return(
					model.User{},
					appErrors.ErrUserNotFound,
				)
				// Create user successfully
				m.On("Create", mock.MatchedBy(func(user model.User) bool {
					return user.Email == "newuser@example.com" && user.Name == "New User"
				})).Return(nil)
			},
			wantError: nil,
		},
		{
			name: "error - user already exists",
			user: model.User{
				Email:    "existing@example.com",
				Password: "password123",
				Name:     "Existing User",
			},
			mockSetup: func(m *MockUserRepositoryAuth) {
				// User already exists
				m.On("GetByEmail", "existing@example.com").Return(
					model.User{
						ID:    1,
						Email: "existing@example.com",
						Name:  "Existing User",
					},
					nil,
				)
			},
			wantError: appErrors.ErrEmailExists,
		},
		{
			name: "error - database error on check",
			user: model.User{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			mockSetup: func(m *MockUserRepositoryAuth) {
				// Database error when checking if user exists
				m.On("GetByEmail", "test@example.com").Return(
					model.User{},
					errors.New("database connection error"),
				)
			},
			wantError: errors.New("database connection error"),
		},
		{
			name: "error - database error on create",
			user: model.User{
				Email:    "newuser@example.com",
				Password: "password123",
				Name:     "New User",
			},
			mockSetup: func(m *MockUserRepositoryAuth) {
				// User doesn't exist
				m.On("GetByEmail", "newuser@example.com").Return(
					model.User{},
					appErrors.ErrUserNotFound,
				)
				// Error creating user
				m.On("Create", mock.MatchedBy(func(user model.User) bool {
					return user.Email == "newuser@example.com"
				})).Return(errors.New("database insert error"))
			},
			wantError: errors.New("database insert error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			service, mockRepo := NewAuthServiceWithMock()
			tt.mockSetup(mockRepo)

			// Execute
			err := service.Register(tt.user)

			// Assert
			if tt.wantError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_Register_PasswordHashing(t *testing.T) {
	// Test that password is properly hashed during registration
	service, mockRepo := NewAuthServiceWithMock()

	originalPassword := "plaintext123"
	user := model.User{
		Email:    "test@example.com",
		Password: originalPassword,
		Name:     "Test User",
	}

	// Mock setup
	mockRepo.On("GetByEmail", "test@example.com").Return(
		model.User{},
		appErrors.ErrUserNotFound,
	)

	var capturedUser model.User
	mockRepo.On("Create", mock.MatchedBy(func(u model.User) bool {
		capturedUser = u
		return true
	})).Return(nil)

	// Execute
	err := service.Register(user)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, originalPassword, capturedUser.Password, "Password should be hashed")
	assert.NotEmpty(t, capturedUser.Password, "Hashed password should not be empty")

	// Verify the hashed password can be compared with original
	err = bcrypt.CompareHashAndPassword([]byte(capturedUser.Password), []byte(originalPassword))
	assert.NoError(t, err, "Hashed password should match original password")

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_JWTSigningError(t *testing.T) {
	// Test JWT signing with valid credentials but check token generation
	service, mockRepo := NewAuthServiceWithMock()

	// Create a hashed password for testing
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Mock setup
	mockRepo.On("GetByEmail", "test@example.com").Return(
		model.User{
			ID:       1,
			Email:    "test@example.com",
			Password: string(hashedPassword),
		},
		nil,
	)

	// Set a valid JWT secret for this test
	originalSecret := os.Getenv("JWTSECRET")
	os.Setenv("JWTSECRET", "test-secret-key")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWTSECRET", originalSecret)
		} else {
			os.Unsetenv("JWTSECRET")
		}
	}()

	// Execute
	token, err := service.Login("test@example.com", "password123")

	// Assert - this should succeed and generate a valid token
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_ValidPassword(t *testing.T) {
	// Test register with a valid long password to ensure bcrypt path is covered
	service, mockRepo := NewAuthServiceWithMock()

	user := model.User{
		Email:    "test@example.com",
		Password: "validpassword123", // Normal password
		Name:     "Test User",
	}

	// Mock setup
	mockRepo.On("GetByEmail", "test@example.com").Return(
		model.User{},
		appErrors.ErrUserNotFound,
	)

	mockRepo.On("Create", mock.MatchedBy(func(u model.User) bool {
		return u.Email == "test@example.com" && u.Name == "Test User"
	})).Return(nil)

	// Execute
	err := service.Register(user)

	// Assert - should succeed
	assert.NoError(t, err)

	// Verify mock expectations
	mockRepo.AssertExpectations(t)
}
