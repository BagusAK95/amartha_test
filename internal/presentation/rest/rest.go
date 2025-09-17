package rest

import (
	loanhttp "github.com/BagusAK95/amarta_test/internal/application/loan/delivery/http"
	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	"github.com/gin-gonic/gin"
)

func NewRouter(loanUsecase loan.ILoanUsecase) *gin.Engine {
	router := gin.Default()

	handler := loanhttp.NewLoanHandler(loanUsecase)

	// API v1 routes
	api := router.Group("/api/v1")
	{
		loans := api.Group("/loan")
		{
			loans.POST("", handler.CreateLoan)
			// loans.GET("/:id", handler.GetLoan)
			// loans.POST("/:id/approve", handler.ApproveLoan)
			// loans.POST("/:id/investments", handler.AddInvestment)
			// loans.POST("/:id/disburse", handler.DisburseLoan)
		}
	}

	return router
}
