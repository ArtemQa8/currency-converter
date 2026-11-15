package api

import (
	"net/http"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/ArtemQa8/currency-converter/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type Handler struct {
	converterService *service.ConverterService
}

func NewHandler(convService *service.ConverterService) *Handler {
	return &Handler{converterService: convService}
}

func (h *Handler) Convert(c *gin.Context) {
	var req ConvertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err.Error())
		return
	}

	amount := decimal.NewFromFloat(req.Amount)
	result, err := h.converterService.Convert(c.Request.Context(), amount, req.From, req.To)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	precision := model.GetPrecision(req.To)
	roundedResult := result.RoundBank(precision)
	c.JSON(http.StatusOK, ConverterResponse{Result: roundedResult.InexactFloat64()})
}

func respondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"ошибка": msg})
}
