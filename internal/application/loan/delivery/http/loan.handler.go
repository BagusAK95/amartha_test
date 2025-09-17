package http

import (
	"net/http"

	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/gin-gonic/gin"
)

type loanHandler struct {
	usecase loan.ILoanUsecase
}

func NewLoanHandler(usecase loan.ILoanUsecase) *loanHandler {
	return &loanHandler{
		usecase: usecase,
	}
}

func (h *loanHandler) CreateLoan(c *gin.Context) {
	var body loan.CreateLoanRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(httpError.NewBadRequestError(err.Error()))
		return
	}

	loan, err := h.usecase.CreateLoan(c.Request.Context(), body)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, loan)
}
