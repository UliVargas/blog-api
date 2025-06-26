package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/UliVargas/blog-go/internal/domain/model"
	"github.com/UliVargas/blog-go/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	require.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

func TestNewUserRepository(t *testing.T) {
	db, _, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewUserRepository(db)

	assert.NotNil(t, repo)
	assert.Equal(t, db, repo.db)
}

func TestUserRepository_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedUsers []model.User
		expectedError error
	}{
		{
			name: "success - returns users",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john@example.com", "password123", time.Now(), time.Now()).
					AddRow(2, "Jane Smith", "jane@example.com", "password456", time.Now(), time.Now())

				mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)
			},
			expectedUsers: []model.User{
				{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "password123"},
				{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Password: "password456"},
			},
			expectedError: nil,
		},
		{
			name: "success - empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnRows(rows)
			},
			expectedUsers: []model.User{},
			expectedError: nil,
		},
		{
			name: "error - database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users"`).WillReturnError(sql.ErrConnDone)
			},
			expectedUsers: nil,
			expectedError: errors.ErrDatabaseOperation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewUserRepository(db)
			tt.setupMock(mock)

			users, err := repo.GetAll()

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.expectedUsers), len(users))
				for i, expectedUser := range tt.expectedUsers {
					assert.Equal(t, expectedUser.ID, users[i].ID)
					assert.Equal(t, expectedUser.Name, users[i].Name)
					assert.Equal(t, expectedUser.Email, users[i].Email)
					assert.Equal(t, expectedUser.Password, users[i].Password)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	tests := []struct {
		name         string
		userID       uint
		setupMock    func(sqlmock.Sqlmock)
		expectedUser model.User
		expectedError error
	}{
		{
			name:   "success - user found",
			userID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john@example.com", "password123", time.Now(), time.Now())

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1`).WillReturnRows(rows)
			},
			expectedUser: model.User{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "password123"},
			expectedError: nil,
		},
		{
			name:   "error - user not found",
			userID: 999,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1`).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser: model.User{},
			expectedError: errors.ErrUserNotFound,
		},
		{
			name:   "error - database error",
			userID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1`).WillReturnError(sql.ErrConnDone)
			},
			expectedUser: model.User{},
			expectedError: errors.ErrDatabaseOperation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewUserRepository(db)
			tt.setupMock(mock)

			user, err := repo.GetByID(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Equal(t, tt.expectedUser, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Password, user.Password)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	tests := []struct {
		name         string
		email        string
		setupMock    func(sqlmock.Sqlmock)
		expectedUser model.User
		expectedError error
	}{
		{
			name:  "success - user found",
			email: "john@example.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john@example.com", "password123", time.Now(), time.Now())

				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).WillReturnRows(rows)
			},
			expectedUser: model.User{ID: 1, Name: "John Doe", Email: "john@example.com", Password: "password123"},
			expectedError: nil,
		},
		{
			name:  "error - user not found",
			email: "notfound@example.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1`).WillReturnError(gorm.ErrRecordNotFound)
			},
			expectedUser: model.User{},
			expectedError: errors.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewUserRepository(db)
			tt.setupMock(mock)

			user, err := repo.GetByEmail(tt.email)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Equal(t, tt.expectedUser, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Password, user.Password)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Create(t *testing.T) {
	tests := []struct {
		name          string
		user          model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name: "success - user created",
			user: model.User{Name: "John Doe", Email: "john@example.com", Password: "password123"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name: "error - database error",
			user: model.User{Name: "John Doe", Email: "existing@example.com", Password: "password123"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedError: errors.ErrDatabaseOperation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewUserRepository(db)
			tt.setupMock(mock)

			err := repo.Create(tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	tests := []struct {
		name          string
		user          model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedUser  model.User
		expectedError error
	}{
		{
			name: "success - user updated",
			user: model.User{ID: 1, Name: "John Updated", Email: "john.updated@example.com", Password: "newpassword"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "users" SET`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedUser: model.User{ID: 1, Name: "John Updated", Email: "john.updated@example.com", Password: "newpassword"},
			expectedError: nil,
		},
		{
			name: "error - database error",
			user: model.User{ID: 1, Name: "John Updated", Email: "john.updated@example.com", Password: "newpassword"},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`UPDATE "users" SET`).WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedUser: model.User{},
			expectedError: errors.ErrDatabaseOperation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewUserRepository(db)
			tt.setupMock(mock)

			user, err := repo.Update(tt.user)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Equal(t, tt.expectedUser, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser.ID, user.ID)
				assert.Equal(t, tt.expectedUser.Name, user.Name)
				assert.Equal(t, tt.expectedUser.Email, user.Email)
				assert.Equal(t, tt.expectedUser.Password, user.Password)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		setupMock     func(sqlmock.Sqlmock)
		expectedError error
	}{
		{
			name:   "success - user deleted",
			userID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "users" WHERE "users"."id" = \$1`).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			expectedError: nil,
		},
		{
			name:   "error - database error",
			userID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM "users" WHERE "users"."id" = \$1`).WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			expectedError: errors.ErrDatabaseOperation,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, cleanup := setupTestDB(t)
			defer cleanup()

			repo := NewUserRepository(db)
			tt.setupMock(mock)

			err := repo.Delete(tt.userID)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}