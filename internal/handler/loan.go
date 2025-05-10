package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github/ramabmtr/billing-engine/internal/lib"
	"github/ramabmtr/billing-engine/internal/model"
	"github/ramabmtr/billing-engine/internal/service"
)

type LoanHandler struct {
	loanSvc *service.LoanService
}

func NewLoanHandler(loanSvc *service.LoanService) *LoanHandler {
	return &LoanHandler{loanSvc: loanSvc}
}

func (h *LoanHandler) RegisterRoutes(g *echo.Group) {
	rg := g.Group("/borrowers/:borrowerID/loans")
	rg.POST("", h.CreateLoanRequest)
	rg.GET("", h.List)
	rg.GET("/:id", h.Detail)
}

// CreateLoanRequest godoc
// @Summary Create a loan request
// @Description Create a new loan request for a specific borrower
// @Tags loans
// @Produce json
// @Param borrowerID path string true "Borrower ID"
// @Success 200 {object} lib.Response "Successfully created loan request"
// @Failure 500 {object} lib.Response "Internal server error"
// @Router /borrowers/{borrowerID}/loans [post]
// @Security ApiKeyAuth
func (h *LoanHandler) CreateLoanRequest(c echo.Context) error {
	borrowerID := c.Param("borrowerID")
	if borrowerID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid borrower ID")
	}
	loan, err := h.loanSvc.CreateLoanRequest(borrowerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, lib.ResponseError(err))
	}

	return c.JSON(http.StatusOK, lib.ResponseSuccess(loan, "loan"))
}

// List godoc
// @Summary List loans for a borrower
// @Description Get a list of all loans for a specific borrower
// @Tags loans
// @Produce json
// @Param borrowerID path string true "Borrower ID"
// @Success 200 {object} lib.Response "Successfully retrieved loans list"
// @Failure 500 {object} lib.Response "Internal server error"
// @Router /borrowers/{borrowerID}/loans [get]
// @Security ApiKeyAuth
func (h *LoanHandler) List(c echo.Context) error {
	borrowerID := c.Param("borrowerID")
	if borrowerID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid borrower ID")
	}
	loans, err := h.loanSvc.GetLoansByBorrowerID(borrowerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, lib.ResponseError(err))
	}

	return c.JSON(http.StatusOK, lib.ResponseSuccess(loans, "loans"))
}

type GetLoanRes struct {
	OutstandingAmount decimal.Decimal `json:"outstanding_amount"`
	Loan              *model.Loan     `json:"loan"`
}

// Detail godoc
// @Summary Get loan details
// @Description Get detailed information about a specific loan
// @Tags loans
// @Produce json
// @Param borrowerID path string true "Borrower ID"
// @Param id path string true "Loan ID"
// @Success 200 {object} lib.Response{data=GetLoanRes} "Successfully retrieved loan details"
// @Failure 500 {object} lib.Response "Internal server error"
// @Router /borrowers/{borrowerID}/loans/{id} [get]
// @Security ApiKeyAuth
func (h *LoanHandler) Detail(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid loan ID")
	}
	loan, outstanding, err := h.loanSvc.GetLoanDetail(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, lib.ResponseError(err))
	}

	return c.JSON(http.StatusOK, lib.ResponseSuccess(GetLoanRes{
		OutstandingAmount: outstanding,
		Loan:              loan,
	}, "loan_detail"))
}
