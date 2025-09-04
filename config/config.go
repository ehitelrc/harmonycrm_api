package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadConfig() {
	dbUser := "usreprac"
	dbPassword := "224wolFe224!" //os.Getenv("DB_PASSWORD") // O ingrésalo directamente en desarrollo
	dbName := "harmony"
	dbHost := "20.118.233.183"
	dbPort := "5432"

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("No se pudo conectar a la base de datos: " + err.Error())
	}

	fmt.Println("Conexión a la base de datos establecida con GORM")
}

// Cierra la conexión utilizando el objeto subyacente sql.DB
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
