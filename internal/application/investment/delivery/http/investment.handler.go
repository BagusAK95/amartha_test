package http

import (
	"net/http"

	"github.com/BagusAK95/amarta_test/internal/domain/investment"
	httpError "github.com/BagusAK95/amarta_test/internal/utils/error"
	"github.com/BagusAK95/amarta_test/internal/utils/html"
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

func (h *investmentHandler) GetInvestmentAgreementFile(c *gin.Context) {
	investmentID, err := uuid.Parse(c.Param("investment_id"))
	if err != nil {
		_ = c.Error(httpError.NewBadRequestError("invalid investment ID"))
		return
	}

	agreementDetail, err := h.usecase.GetInvestmentAgreementDetail(c.Request.Context(), investmentID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	tmpl, err := html.NewTemplate()
	if err != nil {
		_ = c.Error(httpError.NewInternalError("failed to load template"))
		return
	}

	c.Header("Content-Type", "text/html")
	c.Status(http.StatusOK)

	err = tmpl.Execute(c.Writer, "investment_agreement.html", agreementDetail)
	if err != nil {
		_ = c.Error(httpError.NewInternalError("failed to render template"))
		return
	}
}
