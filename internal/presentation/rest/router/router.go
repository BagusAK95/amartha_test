package router

import (
	investmenthttp "github.com/BagusAK95/amarta_test/internal/application/investment/delivery/http"
	loanhttp "github.com/BagusAK95/amarta_test/internal/application/loan/delivery/http"
	"github.com/BagusAK95/amarta_test/internal/domain/investment"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	"github.com/BagusAK95/amarta_test/internal/presentation/rest/middleware"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func NewRouter(loanUsecase loan.ILoanUsecase, investmentUsecase investment.IInvestmentUsecase, tracer opentracing.Tracer) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.TracingMiddleware(tracer))
	router.Use(middleware.ErrorHandler())

	loanHandler := loanhttp.NewLoanHandler(loanUsecase)
	investmentHandler := investmenthttp.NewInvestmentHandler(investmentUsecase)

	// API v1 routes
	api := router.Group("/api/v1")
	{
		loans := api.Group("/loan")
		loans.Use(middleware.AuthMiddleware(middleware.RoleEmployee))
		{
			loans.POST("", loanHandler.CreateLoan)
			loans.GET("", loanHandler.ListLoan)
			loans.GET("/:id", loanHandler.DetailLoan)
			loans.PATCH("/:id/reject", loanHandler.RejectLoan)
			loans.PATCH("/:id/approve", loanHandler.ApproveLoan)
			loans.PATCH("/:id/disburse", loanHandler.DisburseLoan)
		}

		investments := api.Group("/investment")
		investments.Use(middleware.AuthMiddleware(middleware.RoleInvestor))
		{
			investments.POST("", investmentHandler.AddInvestment)
		}

		api.GET("/loan/agreement/file/:loan_id", loanHandler.GetLoanAgreementFile)
		api.GET("/investment/agreement/file/:investment_id", investmentHandler.GetInvestmentAgreementFile)
	}

	return router
}
