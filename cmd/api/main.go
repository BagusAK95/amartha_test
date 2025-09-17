package main

import (
	"fmt"
	"log"

	borrowerrepo "github.com/BagusAK95/amarta_test/internal/application/borrower/repository"
	employeerepo "github.com/BagusAK95/amarta_test/internal/application/employee/repository"
	loanrepo "github.com/BagusAK95/amarta_test/internal/application/loan/repository"
	loanuc "github.com/BagusAK95/amarta_test/internal/application/loan/usecase"
	"github.com/BagusAK95/amarta_test/internal/config"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/database"
	"github.com/BagusAK95/amarta_test/internal/presentation/rest/router"
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
	employeeRepo := employeerepo.NewEmployeeRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	borrowerRepo := borrowerrepo.NewBorrowerRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	loanRepo := loanrepo.NewLoanRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)

	// Initialize usecase
	loanUsecase := loanuc.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)

	// Start server
	gin.SetMode(gin.ReleaseMode)
	r := router.NewRouter(loanUsecase)

	serverAddr := fmt.Sprintf(":%d", cfg.Application.Port)
	log.Printf("🚀 Starting server on %s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
