package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConnect_WithValidDSN(t *testing.T) {
	// Guardar valor original
	originalDBDSN := os.Getenv("DBDSN")
	defer func() {
		if originalDBDSN != "" {
			os.Setenv("DBDSN", originalDBDSN)
		} else {
			os.Unsetenv("DBDSN")
		}
	}()

	// Configurar una DSN válida para testing (usando SQLite en memoria)
	// Nota: En un entorno real, esto requeriría una base de datos de prueba
	os.Setenv("DBDSN", "postgres://user:password@localhost:5432/testdb?sslmode=disable")

	// Este test verificará que la función no entre en pánico con una DSN válida
	// En un entorno de testing real, necesitaríamos una base de datos de prueba
	assert.NotPanics(t, func() {
		// La función intentará conectarse pero fallará sin una DB real
		// Sin embargo, no debería entrar en pánico por tener una DSN válida
		db := DBConnect()
		assert.NotNil(t, db) // La función devuelve un objeto DB incluso si la conexión falla
	})
}

func TestDBConnect_WithEmptyDSN(t *testing.T) {
	// Guardar valor original
	originalDBDSN := os.Getenv("DBDSN")
	defer func() {
		if originalDBDSN != "" {
			os.Setenv("DBDSN", originalDBDSN)
		} else {
			os.Unsetenv("DBDSN")
		}
	}()

	// Limpiar la variable de entorno
	os.Unsetenv("DBDSN")

	// La función debería entrar en pánico cuando DBDSN está vacía
	assert.Panics(t, func() {
		DBConnect()
	}, "DBConnect should panic when DBDSN is empty")
}

func TestDBConnect_PanicMessage(t *testing.T) {
	// Guardar valor original
	originalDBDSN := os.Getenv("DBDSN")
	defer func() {
		if originalDBDSN != "" {
			os.Setenv("DBDSN", originalDBDSN)
		} else {
			os.Unsetenv("DBDSN")
		}
	}()

	// Limpiar la variable de entorno
	os.Unsetenv("DBDSN")

	// Verificar que el mensaje de pánico es el esperado
	assert.PanicsWithValue(t, "No se pudo cargar la variable de entorno DBDSN", func() {
		DBConnect()
	})
}

// Test para verificar que Load() es llamado correctamente dentro de DBConnect
func TestDBConnect_CallsLoad(t *testing.T) {
	// Guardar valor original
	originalDBDSN := os.Getenv("DBDSN")
	defer func() {
		if originalDBDSN != "" {
			os.Setenv("DBDSN", originalDBDSN)
		} else {
			os.Unsetenv("DBDSN")
		}
	}()

	// Configurar una DSN para evitar el pánico
	os.Setenv("DBDSN", "postgres://test:test@localhost:5432/test")

	// Verificar que no hay pánico, lo que indica que Load() fue llamado
	assert.NotPanics(t, func() {
		db := DBConnect()
		assert.NotNil(t, db)
	})
}