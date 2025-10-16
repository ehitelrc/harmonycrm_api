package config

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func LoadConfig() {
	dbUser := "harmony"
	dbPassword := "harmonyValle2025" //os.Getenv("DB_PASSWORD") // O ingrésalo directamente en desarrollo
	dbName := "harmony"
	dbHost := "20.81.232.132"
	dbPort := "5432"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("❌ No se pudo conectar a la base de datos: " + err.Error())
	}

	sqlDB, err := DB.DB()
	if err != nil {
		panic("❌ Error obteniendo instancia sql.DB: " + err.Error())
	}

	// ✅ Configuración del pool de conexiones
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	log.Println("✅ Conexión a PostgreSQL establecida con GORM (pool configurado)")

	// 🔍 Iniciar monitoreo del pool
	go monitorDBConnections(sqlDB)
}

// 🔍 Monitorea periódicamente el estado del pool de conexiones
func monitorDBConnections(sqlDB *sql.DB) {
	for {
		stats := sqlDB.Stats()
		log.Printf("📊 Pool DB — Open: %d | InUse: %d | Idle: %d | WaitCount: %d",
			stats.OpenConnections, stats.InUse, stats.Idle, stats.WaitCount)
		time.Sleep(60 * time.Second)
	}
}

// ❌ Cierra la conexión al detener el servidor
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	log.Println("🧹 Cerrando conexión con PostgreSQL...")
	return sqlDB.Close()
}
