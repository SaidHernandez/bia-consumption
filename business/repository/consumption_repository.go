package repository

import (
	"github.com/SaidHernandez/bia-comsumtion/business/model"
	"github.com/SaidHernandez/bia-comsumtion/infraestructure/db"
)

type ConsumptionRepositoryInterface interface {
	GetConsumptionByFilters(meterID int, startDate, endDate string) ([]model.Consumption, error)
}

type ConsumptionRepository struct{}

func NewConsumptionRepository() *ConsumptionRepository {
	return &ConsumptionRepository{}
}

func (a *ConsumptionRepository) GetConsumptionByFilters(meterID int, startDate, endDate string) ([]model.Consumption, error) {
	var consumptions []model.Consumption
	query := db.DB
	if meterID != 0 {
		query = query.Where("meter_id = ?", meterID)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
	}
	result := query.Find(&consumptions)
	return consumptions, result.Error
}
