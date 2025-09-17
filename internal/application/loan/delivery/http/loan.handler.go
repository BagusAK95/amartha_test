package http

import (
	"net/http"

	"github.com/BagusAK95/amarta_test/internal/domain/loan"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loan, err := h.usecase.CreateLoan(c.Request.Context(), body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, loan)
}
