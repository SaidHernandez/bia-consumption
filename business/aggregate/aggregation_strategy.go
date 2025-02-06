package aggregate

import "github.com/SaidHernandez/bia-comsumtion/business/model"

type AggregationStrategy interface {
	Aggregate(consumptions []model.Consumption) map[string]model.AggregatedConsumption
}
