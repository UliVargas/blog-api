package service

import (
	"errors"
	"testing"

	"github.com/UliVargas/blog-go/internal/domain/model"
	appErrors "github.com/UliVargas/blog-go/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository es un mock para el repositorio de usuarios
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetAll() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) GetByID(id uint) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) Create(user model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// NewUserServiceWithMock creates a UserService with a mock repository for testing
func NewUserServiceWithMock() (*UserService, *MockUserRepository) {
	mockRepo := &MockUserRepository{}
	service := NewUserService(mockRepo)
	return service, mockRepo
}

func TestNewUserService(t *testing.T) {
	userService := &UserService{userRepo: nil}

	assert.NotNil(t, userService)
}

func TestUserService_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		mockSetup     func(*MockUserRepository)
		expectedUsers []model.User
		expectedError error
	}{
		{
			name: "success - users found",
			mockSetup: func(mockRepo *MockUserRepository) {
				users := []model.User{
					{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "hashedpassword"},
					{ID: 2, Name: "Jane Doe", Email: "jane@example.com", Password: "hashedpassword2"},
				}
				mockRepo.On("GetAll").Return(users, nil)
			},
			expectedUsers: []model.User{
				{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "hashedpassword"},
				{ID: 2, Name: "Jane Doe", Email: "jane@example.com", Password: "hashedpassword2"},
			},
			expectedError: nil,
		},
		{
			name: "error - database error",
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetAll").Return([]model.User{}, errors.New("database error"))
			},
			expectedUsers: []model.User{},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService, mockRepo := NewUserServiceWithMock()

			tt.mockSetup(mockRepo)

			users, err := userService.GetAll()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUsers, users)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		mockSetup     func(*MockUserRepository)
		expectedUser  model.User
		expectedError error
	}{
		{
			name:   "success - user found",
			userID: 1,
			mockSetup: func(mockRepo *MockUserRepository) {
				user := model.User{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "hashedpassword"}
				mockRepo.On("GetByID", uint(1)).Return(user, nil)
			},
			expectedUser:  model.User{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "hashedpassword"},
			expectedError: nil,
		},
		{
			name:   "error - user not found",
			userID: 999,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByID", uint(999)).Return(model.User{}, appErrors.ErrUserNotFound)
			},
			expectedUser:  model.User{},
			expectedError: appErrors.ErrUserNotFound,
		},
		{
			name:   "error - database error",
			userID: 1,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("GetByID", uint(1)).Return(model.User{}, errors.New("database error"))
			},
			expectedUser:  model.User{},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService, mockRepo := NewUserServiceWithMock()

			tt.mockSetup(mockRepo)

			user, err := userService.GetByID(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Equal(t, tt.expectedUser, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	tests := []struct {
		name          string
		user          model.User
		mockSetup     func(*MockUserRepository)
		expectedUser  model.User
		expectedError error
	}{
		{
			name: "success - user updated",
			user: model.User{ID: 1, Name: "John Updated", Email: "john.updated@example.com", Password: "newhashedpassword"},
			mockSetup: func(mockRepo *MockUserRepository) {
				updatedUser := model.User{ID: 1, Name: "John Updated", Email: "john.updated@example.com", Password: "newhashedpassword"}
				mockRepo.On("Update", updatedUser).Return(updatedUser, nil)
			},
			expectedUser:  model.User{ID: 1, Name: "John Updated", Email: "john.updated@example.com", Password: "newhashedpassword"},
			expectedError: nil,
		},
		{
			name: "error - database error",
			user: model.User{ID: 1, Name: "John Updated", Email: "john.updated@example.com", Password: "newhashedpassword"},
			mockSetup: func(mockRepo *MockUserRepository) {
				updatedUser := model.User{ID: 1, Name: "John Updated", Email: "john.updated@example.com", Password: "newhashedpassword"}
				mockRepo.On("Update", updatedUser).Return(model.User{}, errors.New("database error"))
			},
			expectedUser:  model.User{},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService, mockRepo := NewUserServiceWithMock()

			tt.mockSetup(mockRepo)

			user, err := userService.Update(tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Equal(t, tt.expectedUser, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		mockSetup     func(*MockUserRepository)
		expectedError error
	}{
		{
			name:   "success - user deleted",
			userID: 1,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("Delete", uint(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "error - database error",
			userID: 1,
			mockSetup: func(mockRepo *MockUserRepository) {
				mockRepo.On("Delete", uint(1)).Return(errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService, mockRepo := NewUserServiceWithMock()

			tt.mockSetup(mockRepo)

			err := userService.Delete(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
