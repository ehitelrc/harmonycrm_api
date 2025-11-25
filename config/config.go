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

// LoadConfig inicializa la conexión global a PostgreSQL con GORM
// y configura un pool estable para alta concurrencia.
func LoadConfig() {
	// ⚠️ En producción, ideal usar variables de entorno
	dbUser := "harmony"
	dbPassword := "harmonyValle2025" // os.Getenv("DB_PASSWORD")
	dbName := "harmony"
	dbHost := "20.81.232.132"
	dbPort := "5432"

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPassword, dbName, dbPort,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Opcional: mejora performance en alto tráfico
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a la base de datos: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("❌ Error obteniendo instancia sql.DB: %v", err)
	}

	// ✅ Configuración del pool para servicio crítico

	// Máximo de conexiones abiertas hacia PostgreSQL.
	// IMPORTANTE: no debe superar (max_connections - conexiones de otros servicios).
	// Ajustá según tu servidor (ej: si max_connections=200, dejar 80–100 aquí).
	sqlDB.SetMaxOpenConns(80)

	// Máximo de conexiones inactivas que se mantienen en el pool.
	sqlDB.SetMaxIdleConns(20)

	// Tiempo máximo de vida de una conexión (rotación suave).
	// Evita problemas con proxies/NAT de larga data sin matar conexiones muy seguido.
	sqlDB.SetConnMaxLifetime(55 * time.Minute)

	// Tiempo máximo que una conexión puede permanecer idle antes de cerrarse.
	sqlDB.SetConnMaxIdleTime(15 * time.Minute)

	log.Println("✅ Conexión a PostgreSQL establecida con GORM (pool configurado para alta concurrencia)")

	// 🔍 Iniciar monitoreo del pool + ping periódico
	go monitorDBConnections(sqlDB)
	go healthCheckLoop(sqlDB)
}

// 🔍 Monitorea periódicamente el estado del pool de conexiones
func monitorDBConnections(sqlDB *sql.DB) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		stats := sqlDB.Stats()
		log.Printf(
			"📊 Pool DB — Open=%d | InUse=%d | Idle=%d | WaitCount=%d | WaitDuration=%s",
			stats.OpenConnections,
			stats.InUse,
			stats.Idle,
			stats.WaitCount,
			stats.WaitDuration,
		)
	}
}

// 🩺 Verifica periódicamente que la DB siga respondiendo
// Útil para detectar que el pool se rompió o el servidor DB cayó.
func healthCheckLoop(sqlDB *sql.DB) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := sqlDB.Ping(); err != nil {
			log.Printf("⚠️ HealthCheck DB: error al hacer ping a la base de datos: %v", err)
			// Aquí NO reabrimos la conexión para no generar condiciones de carrera.
			// Lo ideal es que un orquestador (Docker/systemd/K8s) reinicie el servicio
			// si detecta demasiados errores.
		}
	}
}

// ❌ Cierra la conexión al detener el servidor
func CloseDB() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	log.Println("🧹 Cerrando conexión con PostgreSQL...")
	return sqlDB.Close()
}
