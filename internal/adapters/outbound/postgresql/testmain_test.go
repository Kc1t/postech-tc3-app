package postgresql

import (
	"fmt"
	"log"
	"os"
	"testing"

	pgmodel "github.com/fiap/postech-tc1/internal/adapters/outbound/postgresql/model"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	host := getEnvOrDefault("TEST_DB_HOST", "localhost")
	port := getEnvOrDefault("TEST_DB_PORT", "5432")
	user := getEnvOrDefault("TEST_DB_USER", "test")
	password := getEnvOrDefault("TEST_DB_PASSWORD", "test")
	dbname := getEnvOrDefault("TEST_DB_NAME", "testdb")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("failed to connect to test database: %v", err)
	}

	if err := db.AutoMigrate(
		&pgmodel.Customer{},
		&pgmodel.Vehicle{},
		&pgmodel.Part{},
		&pgmodel.Service{},
		&pgmodel.ServiceOrder{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	testDB = db

	os.Exit(m.Run())
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
