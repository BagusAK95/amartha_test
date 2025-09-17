package database

import (
	"github.com/BagusAK95/amarta_test/internal/config"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/database/postgres"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Postgres postgres.PostgresConfig
}

type DatabaseConnection struct {
	Postgres postgres.PostgresConnection
}

func SetConfig(cfg config.PostgresConfig) DatabaseConfig {
	return DatabaseConfig{
		Postgres: postgres.SetConfig(cfg),
	}
}

func OpenConnection(cfg config.PostgresConfig, dbConfig DatabaseConfig) DatabaseConnection {
	return DatabaseConnection{
		Postgres: postgres.PostgresConnection{
			Master: postgres.Open(cfg, dbConfig.Postgres.Master),
			Slave:  postgres.Open(cfg, dbConfig.Postgres.Slave),
		},
	}
}

func CloseConnection(dbConn DatabaseConnection) {
	if dbConn.Postgres.Master != nil {
		postgres.Close(dbConn.Postgres.Master)
	}

	if dbConn.Postgres.Slave != nil {
		postgres.Close(dbConn.Postgres.Slave)
	}
}
