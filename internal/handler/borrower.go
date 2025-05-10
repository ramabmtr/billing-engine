package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github/ramabmtr/billing-engine/internal/lib"
	"github/ramabmtr/billing-engine/internal/service"
)

type BorrowerHandler struct {
	borrowerSvc *service.BorrowerService
}

func NewBorrowerHandler(borrowerSvc *service.BorrowerService) *BorrowerHandler {
	return &BorrowerHandler{borrowerSvc: borrowerSvc}
}

func (h *BorrowerHandler) RegisterRoutes(g *echo.Group) {
	rg := g.Group("/borrowers")
	rg.POST("", h.Create)
	rg.GET("", h.List)
}

type CreateBorrowerReqBody struct {
	Name string `json:"name" validate:"required"`
}

// Create godoc
// @Summary Create a new borrower
// @Description Create a new borrower with the provided name
// @Tags borrowers
// @Accept json
// @Produce json
// @Param request body CreateBorrowerReqBody true "Borrower information"
// @Success 200 {object} lib.Response "Successfully created borrower"
// @Failure 500 {object} lib.Response "Internal server error"
// @Router /borrowers [post]
func (h *BorrowerHandler) Create(c echo.Context) error {
	var req CreateBorrowerReqBody
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}
	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	borrower, err := h.borrowerSvc.Create(req.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, lib.ResponseError(err))
	}

	return c.JSON(http.StatusOK, lib.ResponseSuccess(borrower, "borrower"))
}

// List godoc
// @Summary List all borrowers
// @Description Get a list of all borrowers
// @Tags borrowers
// @Produce json
// @Success 200 {object} lib.Response "Successfully retrieved borrowers list"
// @Failure 500 {object} lib.Response "Internal server error"
// @Router /borrowers [get]
func (h *BorrowerHandler) List(c echo.Context) error {
	borrowers, err := h.borrowerSvc.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, lib.ResponseError(err))
	}

	return c.JSON(http.StatusOK, lib.ResponseSuccess(borrowers, "borrowers"))
}
