package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	// Guardar valores originales de las variables de entorno
	originalDBDSN := os.Getenv("DBDSN")
	originalJWTSECRET := os.Getenv("JWTSECRET")
	originalPORT := os.Getenv("PORT")

	// Limpiar al final del test
	defer func() {
		if originalDBDSN != "" {
			os.Setenv("DBDSN", originalDBDSN)
		} else {
			os.Unsetenv("DBDSN")
		}
		if originalJWTSECRET != "" {
			os.Setenv("JWTSECRET", originalJWTSECRET)
		} else {
			os.Unsetenv("JWTSECRET")
		}
		if originalPORT != "" {
			os.Setenv("PORT", originalPORT)
		} else {
			os.Unsetenv("PORT")
		}
	}()

	t.Run("load config with all environment variables set", func(t *testing.T) {
		// Configurar variables de entorno para el test
		os.Setenv("DBDSN", "postgres://user:password@localhost/testdb")
		os.Setenv("JWTSECRET", "test-jwt-secret")
		os.Setenv("PORT", "8080")

		config := Load()

		assert.NotNil(t, config)
		assert.Equal(t, "postgres://user:password@localhost/testdb", config.DBDSN)
		assert.Equal(t, "test-jwt-secret", config.JWTSECRET)
		assert.Equal(t, "8080", config.PORT)
	})

	t.Run("load config with empty environment variables", func(t *testing.T) {
		// Limpiar variables de entorno
		os.Unsetenv("DBDSN")
		os.Unsetenv("JWTSECRET")
		os.Unsetenv("PORT")

		config := Load()

		assert.NotNil(t, config)
		assert.Equal(t, "", config.DBDSN)
		assert.Equal(t, "", config.JWTSECRET)
		assert.Equal(t, "", config.PORT)
	})

	t.Run("load config with partial environment variables", func(t *testing.T) {
		// Configurar solo algunas variables
		os.Setenv("DBDSN", "postgres://localhost/db")
		os.Unsetenv("JWTSECRET")
		os.Setenv("PORT", "3000")

		config := Load()

		assert.NotNil(t, config)
		assert.Equal(t, "postgres://localhost/db", config.DBDSN)
		assert.Equal(t, "", config.JWTSECRET)
		assert.Equal(t, "3000", config.PORT)
	})
}

func TestConfig_Structure(t *testing.T) {
	config := &Config{
		DBDSN:     "test-dsn",
		JWTSECRET: "test-secret",
		PORT:      "8080",
	}

	assert.Equal(t, "test-dsn", config.DBDSN)
	assert.Equal(t, "test-secret", config.JWTSECRET)
	assert.Equal(t, "8080", config.PORT)
	assert.IsType(t, &Config{}, config)
}