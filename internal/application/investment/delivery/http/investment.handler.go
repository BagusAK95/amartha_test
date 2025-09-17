package http

import (
	"net/http"

	"github.com/BagusAK95/amarta_test/internal/domain/investment"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type investmentHandler struct {
	usecase investment.IInvestmentUsecase
}

func NewInvestmentHandler(usecase investment.IInvestmentUsecase) *investmentHandler {
	return &investmentHandler{
		usecase: usecase,
	}
}

func (h *investmentHandler) AddInvestment(c *gin.Context) {
	var body investment.CreateInvestmentRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		_ = c.Error(httpError.NewBadRequestError(err.Error()))
		return
	}

	investorID, _ := c.Get("investorID")

	res, err := h.usecase.AddInvestment(c.Request.Context(), investorID.(uuid.UUID), body)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, res)
}
