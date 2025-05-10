package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ramabmtr/billing-engine/internal/lib"
	"github.com/ramabmtr/billing-engine/internal/service"
	"github.com/shopspring/decimal"
)

type PaymentHandler struct {
	loanSvc *service.LoanService
}

func NewPaymentHandler(loanSvc *service.LoanService) *PaymentHandler {
	return &PaymentHandler{loanSvc: loanSvc}
}

func (h *PaymentHandler) RegisterRoutes(g *echo.Group) {
	rg := g.Group("/borrowers/:borrowerID/loans/:loanID/payments")
	rg.POST("", h.MakePayment)
	rg.GET("", h.List)
}

type MakePaymentReqBody struct {
	Amount float64 `json:"amount" validate:"required"`
}

// MakePayment godoc
// @Summary Make a payment for a loan
// @Description Process a payment for a specific loan
// @Tags payments
// @Accept json
// @Produce json
// @Param borrowerID path string true "Borrower ID"
// @Param loanID path string true "Loan ID"
// @Param request body MakePaymentReqBody true "Payment information"
// @Success 200 {object} lib.Response "Successfully processed payment"
// @Failure 500 {object} lib.Response "Internal server error"
// @Router /borrowers/{borrowerID}/loans/{loanID}/payments [post]
// @Security ApiKeyAuth
func (h *PaymentHandler) MakePayment(c echo.Context) error {
	loanID := c.Param("loanID")
	if loanID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid loan ID")
	}
	var req MakePaymentReqBody
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := h.loanSvc.MakePayment(loanID, decimal.NewFromFloat(req.Amount))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, lib.ResponseError(err))
	}

	return c.JSON(http.StatusOK, lib.ResponseSuccess(nil))
}

// List godoc
// @Summary List payments for a loan
// @Description Get a list of all payments for a specific loan
// @Tags payments
// @Produce json
// @Param borrowerID path string true "Borrower ID"
// @Param loanID path string true "Loan ID"
// @Success 200 {object} lib.Response "Successfully retrieved payments list"
// @Failure 500 {object} lib.Response "Internal server error"
// @Router /borrowers/{borrowerID}/loans/{loanID}/payments [get]
// @Security ApiKeyAuth
func (h *PaymentHandler) List(c echo.Context) error {
	loanID := c.Param("loanID")
	if loanID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid loan ID")
	}
	payments, err := h.loanSvc.GetLoanPaymentsByLoanID(loanID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, lib.ResponseError(err))
	}

	return c.JSON(http.StatusOK, lib.ResponseSuccess(payments, "payments"))
}
