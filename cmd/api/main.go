package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	borrowerrepo "github.com/BagusAK95/amarta_test/internal/application/borrower/repository"
	employeerepo "github.com/BagusAK95/amarta_test/internal/application/employee/repository"
	investmentrepo "github.com/BagusAK95/amarta_test/internal/application/investment/repository"
	investmentuc "github.com/BagusAK95/amarta_test/internal/application/investment/usecase"
	investorrepo "github.com/BagusAK95/amarta_test/internal/application/investor/repository"
	loanrepo "github.com/BagusAK95/amarta_test/internal/application/loan/repository"
	loanuc "github.com/BagusAK95/amarta_test/internal/application/loan/usecase"
	mailuc "github.com/BagusAK95/amarta_test/internal/application/mail/usecase"
	"github.com/BagusAK95/amarta_test/internal/config"
	"github.com/BagusAK95/amarta_test/internal/domain/mail"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/bus"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/database"
	mailsender "github.com/BagusAK95/amarta_test/internal/infrastructure/mail"
	"github.com/BagusAK95/amarta_test/internal/infrastructure/tracer"
	buslistener "github.com/BagusAK95/amarta_test/internal/presentation/messaging/bus"
	"github.com/BagusAK95/amarta_test/internal/presentation/rest/router"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
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

	// Jaeger
	tracer, closer := tracer.Init(cfg.Jaeger)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// Mail server
	mailSender := mailsender.NewSender(cfg.Mail)
	mailBus := bus.NewBus[mail.MailSendRequest]()

	// Initialize repository
	employeeRepo := employeerepo.NewEmployeeRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	borrowerRepo := borrowerrepo.NewBorrowerRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	loanRepo := loanrepo.NewLoanRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	investmentRepo := investmentrepo.NewInvestmentRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)
	investorRepo := investorrepo.NewInvestorRepo(dbConn.Postgres.Master, dbConn.Postgres.Slave)

	// Initialize usecase
	loanUsecase := loanuc.NewLoanUsecase(loanRepo, borrowerRepo, employeeRepo)
	investmentUsecase := investmentuc.NewInvestmentUsecase(investmentRepo, investorRepo, loanRepo, borrowerRepo, mailBus)
	mailUsecase := mailuc.NewMailUsecase(mailSender)

	// Bus listener
	buslistener.NewBusListener(mailBus, mailUsecase)

	// Start server
	gin.SetMode(gin.ReleaseMode)

	r := router.NewRouter(loanUsecase, investmentUsecase, tracer)
	addr := fmt.Sprintf(":%d", cfg.Application.Port)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Starting server on %s\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("💤 Shutting down server...")

	// Close connection
	database.CloseConnection(dbConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exiting")
}
