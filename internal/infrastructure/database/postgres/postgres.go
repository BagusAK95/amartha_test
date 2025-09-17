package postgres

import (
	"fmt"
	"log"
	"os"

	"github.com/BagusAK95/amarta_test/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

type PostgresConfig struct {
	Master postgres.Config
	Slave  postgres.Config
}

type PostgresConnection struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func SetConfig(cfg config.PostgresConfig) PostgresConfig {
	dsn := "host=%s user=%s password=%s dbname=%s port=%d sslmode=%s timezone=%s"

	return PostgresConfig{
		Master: postgres.Config{
			DSN:                  fmt.Sprintf(dsn, cfg.MasterHost, cfg.MasterUsername, cfg.MasterPassword, cfg.Database, cfg.MasterPort, cfg.MasterSSLMode, cfg.Timezone),
			PreferSimpleProtocol: true,
		},
		Slave: postgres.Config{
			DSN:                  fmt.Sprintf(dsn, cfg.SlaveHost, cfg.SlaveUsername, cfg.SlavePassword, cfg.Database, cfg.SlavePort, cfg.SlaveSSLMode, cfg.Timezone),
			PreferSimpleProtocol: true,
		},
	}
}

func Open(cfg config.PostgresConfig, pgConfig postgres.Config) *gorm.DB {
	db, err := gorm.Open(postgres.New(pgConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("gorm.Open() got error - postgres.Open()", err)
		os.Exit(1)
	}

	if err := db.Use(tracing.NewPlugin(tracing.WithoutQueryVariables())); err != nil {
		log.Fatal("db.Use() got error - postgres.Open()", err)
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("db.DB() got error - postgres.Open()", err)
		os.Exit(1)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	if err = sqlDB.Ping(); err != nil {
		log.Fatal("sqlDB.Ping() got error - postgres.Open()", err)
		os.Exit(1)
	}

	return db
}

func Close(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatal("conn.DB() got error - postgres.Close()", err)
		os.Exit(1)
	}

	sqlDB.Close()
}
