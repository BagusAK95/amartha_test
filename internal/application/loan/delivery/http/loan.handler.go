package http

import (
	"net/http"

	"github.com/BagusAK95/amarta_test/internal/domain/loan"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	res, err := h.usecase.CreateLoan(c.Request.Context(), body)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

func (h *loanHandler) RejectLoan(c *gin.Context) {
	loanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(httpError.NewBadRequestError(err.Error()))
		return
	}

	var body loan.RejectLoanRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(httpError.NewBadRequestError(err.Error()))
		return
	}

	res, err := h.usecase.RejectLoan(c.Request.Context(), loanID, body.RejectReason)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *loanHandler) ApproveLoan(c *gin.Context) {
	loanID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.Error(httpError.NewBadRequestError(err.Error()))
		return
	}

	var body loan.ApproveLoanRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(httpError.NewBadRequestError(err.Error()))
		return
	}

	res, err := h.usecase.ApproveLoan(c.Request.Context(), loanID, body)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, res)
}
