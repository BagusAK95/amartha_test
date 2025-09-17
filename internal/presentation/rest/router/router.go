package router

import (
	loanhttp "github.com/BagusAK95/amarta_test/internal/application/loan/delivery/http"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	"github.com/BagusAK95/amarta_test/internal/presentation/rest/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(loanUsecase loan.ILoanUsecase) *gin.Engine {
	router := gin.Default()
	router.Use(middleware.ErrorHandler())

	handler := loanhttp.NewLoanHandler(loanUsecase)

	// API v1 routes
	api := router.Group("/api/v1")
	{
		loans := api.Group("/loan")
		{
			loans.POST("", handler.CreateLoan)
			loans.GET("", handler.ListLoan)
			loans.GET("/:id", handler.DetailLoan)
			loans.PATCH("/:id/reject", handler.RejectLoan)
			loans.PATCH("/:id/approve", handler.ApproveLoan)
			// loans.POST("/:id/investments", handler.AddInvestment)
			// loans.POST("/:id/disburse", handler.DisburseLoan)
		}
	}

	return router
}
