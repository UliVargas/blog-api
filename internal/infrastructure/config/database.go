package config

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnect() *gorm.DB {
	cfg := Load()
	if cfg.DBDSN == "" {
		panic("No se pudo cargar la variable de entorno DBDSN")
	}

	db, err := gorm.Open(postgres.Open(cfg.DBDSN), &gorm.Config{})
	if err != nil {
		log.Println("No se pudo conectar a la base de datos", err)
	}

	return db
}