package main

import (
	"fmt"
	"log"

	borrowerrepo "github.com/BagusAK95/amarta_test/internal/application/borrower/repository"
	employeerepo "github.com/BagusAK95/amarta_test/internal/application/employee/repository"
	investmentrepo "github.com/BagusAK95/amarta_test/internal/application/investment/repository"
	investmentuc "github.com/BagusAK95/amarta_test/internal/application/investment/usecase"
	investorrepo "github.com/BagusAK95/amarta_test/internal/application/investor/repository"
	loanrepo "github.com/BagusAK95/amarta_test/internal/application/loan/repository"
	loanuc "github.com/BagusAK95/amarta_test/internal/application/loan/usecase"
	"github.com/BagusAK95/amarta_test/internal/config"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/database"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/mail"
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

	mailSender := mail.NewSender(cfg.Mail)

	// Initialize repository
	employeeRepo := employeerepo.NewEmployeeRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	borrowerRepo := borrowerrepo.NewBorrowerRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	loanRepo := loanrepo.NewLoanRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	investmentRepo := investmentrepo.NewInvestmentRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	investorRepo := investorrepo.NewInvestorRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)

	// Initialize usecase
	loanUsecase := loanuc.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
	investmentUsecase := investmentuc.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailSender)

	// Start server
	gin.SetMode(gin.ReleaseMode)
	r := router.NewRouter(loanUsecase, investmentUsecase)

	serverAddr := fmt.Sprintf(":%d", cfg.Application.Port)
	log.Printf("🚀 Starting server on %s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
