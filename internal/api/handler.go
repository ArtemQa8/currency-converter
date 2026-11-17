package api

//nolint:depguard
import (
	"net/http"
	"strconv"

	"github.com/ArtemQa8/currency-converter/internal/model"
	"github.com/ArtemQa8/currency-converter/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// Handler описывает HTTP-Хендлеры.
type Handler struct {
	converterService *service.ConverterService
}

// NewHandler конструктор - создает новый хендлер.
func NewHandler(convService *service.ConverterService) *Handler {
	return &Handler{converterService: convService}
}

// Convert обрабатывает запрос на конвертацию.
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

// History обрабатывает запрос на получение истории курсов.
func (h *Handler) History(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	daysStr := c.DefaultQuery("days", "7")

	if from == "" || to == "" {
		respondError(c, http.StatusBadRequest, "параметры 'from' и 'to' обязательны")
		return
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		respondError(c, http.StatusBadRequest, "параметр 'days' должен быть больше нуля")
		return
	}

	rates, err := h.converterService.GetHistory(c.Request.Context(), from, to, days)
	if err != nil {
		respondError(c, http.StatusInternalServerError, err.Error())
		return
	}

	var response HistoryResponse
	for _, r := range rates {
		precision := model.GetPrecision(r.QuoteCurrency)
		roundedValue := r.Value.RoundBank(precision)
		response.Rates = append(response.Rates, HistoryObject{
			Date:  r.Date,
			From:  r.BaseCurrency,
			To:    r.QuoteCurrency,
			Value: roundedValue.InexactFloat64(),
		})
	}

	c.JSON(http.StatusOK, response)
}
