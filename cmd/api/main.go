package main

import (
	"fmt"
	"log"

	loanrepo "github.com/BagusAK95/amarta_test/internal/application/loan/repository"
	loanuc "github.com/BagusAK95/amarta_test/internal/application/loan/usecase"
	"github.com/BagusAK95/amarta_test/internal/config"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/database"
	"github.com/BagusAK95/amarta_test/internal/presentation/rest"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Could not load config: %v", err)
	}

	// Open database connection
	dbConfig := database.SetConfig(cfg.Postgres)
	dbConn := database.OpenConnection(cfg.Postgres, dbConfig)

	// Initialize repository
	loanRepo := loanrepo.NewLoanRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)

	// Initialize usecase
	loanUsecase := loanuc.NewLoanUsecase(loanRepo)

	// Start server
	gin.SetMode(gin.ReleaseMode)
	r := rest.NewRouter(loanUsecase)

	serverAddr := fmt.Sprintf(":%d", cfg.Application.Port)
	log.Printf("🚀 Starting server on %s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
