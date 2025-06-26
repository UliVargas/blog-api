package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)

	// Guardar valor original de JWT_SECRET
	originalJWTSecret := os.Getenv("JWTSECRET")
	defer func() {
		if originalJWTSecret != "" {
			os.Setenv("JWTSECRET", originalJWTSecret)
		} else {
			os.Unsetenv("JWTSECRET")
		}
	}()

	// Configurar JWT_SECRET para los tests
	testSecret := "test-jwt-secret-key"
	os.Setenv("JWTSECRET", testSecret)

	// Helper function para crear un token JWT válido
	createValidToken := func(userID float64) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour).Unix(),
		})
		tokenString, _ := token.SignedString([]byte(testSecret))
		return tokenString
	}

	// Helper function para crear un token JWT expirado
	createExpiredToken := func(userID float64) string {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(-time.Hour).Unix(), // Expirado hace una hora
		})
		tokenString, _ := token.SignedString([]byte(testSecret))
		return tokenString
	}

	// Helper function para crear un token con algoritmo incorrecto
	createInvalidAlgorithmToken := func(userID float64) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"user_id": userID,
			"exp":     time.Now().Add(time.Hour).Unix(),
		})
		// Esto fallará porque estamos usando RS256 pero el middleware espera HMAC
		tokenString, _ := token.SignedString([]byte(testSecret))
		return tokenString
	}

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		setupEnv       func()
		expectedStatus int
		expectedBody   string
		checkUserID    bool
		expectedUserID uint
	}{
		{
			name: "valid token",
			setupRequest: func(req *http.Request) {
				token := createValidToken(123)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			setupEnv:       func() {},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"success"}`,
			checkUserID:    true,
			expectedUserID: 123,
		},
		{
			name: "missing authorization header",
			setupRequest: func(req *http.Request) {
				// No configurar Authorization header
			},
			setupEnv:       func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Token requerido"}`,
			checkUserID:    false,
		},
		{
			name: "invalid token format - no Bearer",
			setupRequest: func(req *http.Request) {
				token := createValidToken(123)
				req.Header.Set("Authorization", token) // Sin "Bearer "
			},
			setupEnv:       func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Formato de token invalido"}`,
			checkUserID:    false,
		},
		{
			name: "invalid token format - wrong prefix",
			setupRequest: func(req *http.Request) {
				token := createValidToken(123)
				req.Header.Set("Authorization", "Basic "+token)
			},
			setupEnv:       func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Formato de token invalido"}`,
			checkUserID:    false,
		},
		{
			name: "missing JWT secret",
			setupRequest: func(req *http.Request) {
				token := createValidToken(123)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			setupEnv: func() {
				os.Unsetenv("JWTSECRET")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"No se pudo iniciar sesión"}`,
			checkUserID:    false,
		},
		{
			name: "invalid token - malformed",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid.token.here")
			},
			setupEnv:       func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Token inválido"}`,
			checkUserID:    false,
		},
		{
			name: "expired token",
			setupRequest: func(req *http.Request) {
				token := createExpiredToken(123)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			setupEnv:       func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Token inválido"}`,
			checkUserID:    false,
		},
		{
			name: "token with wrong signing method",
			setupRequest: func(req *http.Request) {
				token := createInvalidAlgorithmToken(123)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			setupEnv:       func() {},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Token inválido"}`,
			checkUserID:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Restaurar el entorno después de cada test
			os.Setenv("JWTSECRET", testSecret)
			tt.setupEnv()

			// Crear router y middleware
			router := gin.New()
			router.Use(AuthMiddleware())

			// Endpoint de prueba
			router.GET("/test", func(c *gin.Context) {
				if tt.checkUserID {
					userID, exists := c.Get("user_id")
					assert.True(t, exists, "user_id should be set in context")
					assert.Equal(t, tt.expectedUserID, userID)
				}
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Crear request
			req := httptest.NewRequest("GET", "/test", nil)
			tt.setupRequest(req)

			// Crear response recorder
			w := httptest.NewRecorder()

			// Ejecutar request
			router.ServeHTTP(w, req)

			// Verificar respuesta
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus != http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestAuthMiddleware_TokenWithoutUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Configurar JWT_SECRET
	testSecret := "test-jwt-secret-key"
	os.Setenv("JWTSECRET", testSecret)
	defer os.Unsetenv("JWTSECRET")

	// Crear token sin user_id claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"some_other_claim": "value",
		"exp":              time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(testSecret))

	// Crear router y middleware
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		// Verificar que user_id no está en el contexto
		userID, exists := c.Get("user_id")
		assert.False(t, exists, "user_id should not be set when not in token")
		assert.Nil(t, userID)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Crear request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Crear response recorder
	w := httptest.NewRecorder()

	// Ejecutar request
	router.ServeHTTP(w, req)

	// Verificar que la request pasa pero sin user_id
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMiddleware_TokenWithNonMapClaims(t *testing.T) {
	// Configurar Gin en modo test
	gin.SetMode(gin.TestMode)

	// Configurar JWT_SECRET para el test
	testSecret := "test-jwt-secret-key"
	os.Setenv("JWTSECRET", testSecret)
	defer os.Unsetenv("JWTSECRET")

	// Crear un token con claims personalizados (no MapClaims)
	type CustomClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}

	claims := CustomClaims{
		UserID: 123,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testSecret))

	// Crear router y middleware
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		// El middleware debería funcionar normalmente con CustomClaims
		// porque jwt.Parse convierte automáticamente a MapClaims
		userID, exists := c.Get("user_id")
		if exists {
			// Si existe, debería ser el valor correcto (puede ser uint o float64)
			switch v := userID.(type) {
			case float64:
				assert.Equal(t, float64(123), v)
			case uint:
				assert.Equal(t, uint(123), v)
			default:
				t.Errorf("Unexpected type for user_id: %T", v)
			}
		}
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Crear request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Crear response recorder
	w := httptest.NewRecorder()

	// Ejecutar request
	router.ServeHTTP(w, req)

	// Verificar que la request pasa pero sin user_id
	assert.Equal(t, http.StatusOK, w.Code)
}



func TestAuthMiddleware_TokenWithInvalidUserIDType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Configurar JWT_SECRET
	testSecret := "test-jwt-secret-key"
	os.Setenv("JWTSECRET", testSecret)
	defer os.Unsetenv("JWTSECRET")

	// Crear token con user_id como string en lugar de float64
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "123", // String en lugar de número
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(testSecret))

	// Crear router y middleware
	router := gin.New()
	router.Use(AuthMiddleware())
	router.GET("/test", func(c *gin.Context) {
		// Verificar que user_id no está en el contexto porque no se pudo convertir
		userID, exists := c.Get("user_id")
		assert.False(t, exists, "user_id should not be set when type conversion fails")
		assert.Nil(t, userID)
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Crear request
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)

	// Crear response recorder
	w := httptest.NewRecorder()

	// Ejecutar request
	router.ServeHTTP(w, req)

	// Verificar que la request pasa pero sin user_id
	assert.Equal(t, http.StatusOK, w.Code)
}