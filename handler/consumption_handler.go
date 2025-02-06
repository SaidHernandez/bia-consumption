package handlers

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/SaidHernandez/bia-comsumtion/services"
	"github.com/labstack/echo/v4"
)

// ConsumptionHandler maneja las solicitudes relacionadas con el consumo de energía.
type ConsumptionHandler struct {
	service *services.ConsumptionService
}

// NewConsumptionHandler crea una nueva instancia de ConsumptionHandler.
func NewConsumptionHandler(service *services.ConsumptionService) *ConsumptionHandler {
	return &ConsumptionHandler{service: service}
}

// GetConsumption maneja la solicitud para obtener el consumo por periodo.
// @Summary Obtiene el consumo de energía por periodo.
// @Description Retorna el consumo de energía de los medidores en el rango de fechas especificado.
// @Tags consumption
// @Accept json
// @Produce json
// @Param meter_ids query string true "IDs de los medidores separados por comas"
// @Param start_date query string true "Fecha de inicio en formato YYYY-MM-DD"
// @Param end_date query string true "Fecha de fin en formato YYYY-MM-DD"
// @Param kind_period query string true "Tipo de periodo: daily, weekly, monthly"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /consumption [get]
func (h *ConsumptionHandler) GetConsumption(c echo.Context) error {
	ctx := context.Background()

	meterIDsStr := c.QueryParam("meters_ids")
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	kindPeriod := c.QueryParam("kind_period")

	if meterIDsStr == "" || startDate == "" || endDate == "" || kindPeriod == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Todos los parámetros son requeridos"})
	}
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Formato inválido de start_date, debe ser YYYY-MM-DD"})
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Formato inválido de end_date, debe ser YYYY-MM-DD"})
	}

	if start.After(end) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "start_date no puede ser mayor que end_date"})
	}

	var meterIDs []int
	meterIDsList := strings.Split(meterIDsStr, ",")
	for _, idStr := range meterIDsList {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Formato inválido de meter_ids"})
		}
		meterIDs = append(meterIDs, id)
	}

	results, err := h.service.GetConsumptionByPeriod(ctx, meterIDs, startDate, endDate, kindPeriod)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, results)
}
