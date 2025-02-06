package aggregate

import (
	"github.com/SaidHernandez/bia-comsumtion/business/model"
)

type MonthlyAggregationStrategy struct{}

func (m *MonthlyAggregationStrategy) Aggregate(consumptions []model.Consumption) map[string]model.AggregatedConsumption {
	aggregation := make(map[string]model.AggregatedConsumption)

	for _, consumption := range consumptions {
		month := consumption.Date.Format("Jan 2006")

		if _, exists := aggregation[month]; !exists {
			aggregation[month] = model.AggregatedConsumption{
				Period:             []string{month},
				ActiveEnergy:       []float64{consumption.ActiveEnergy},
				ReactiveInductive:  []float64{consumption.ReactiveInductive},
				ReactiveCapacitive: []float64{consumption.ReactiveCapacitive},
				ExportedEnergy:     []float64{consumption.ExportedEnergy},
			}
		} else {
			aggData := aggregation[month]
			aggData.ActiveEnergy = append(aggData.ActiveEnergy, consumption.ActiveEnergy)
			aggData.ReactiveInductive = append(aggData.ReactiveInductive, consumption.ReactiveInductive)
			aggData.ReactiveCapacitive = append(aggData.ReactiveCapacitive, consumption.ReactiveCapacitive)
			aggData.ExportedEnergy = append(aggData.ExportedEnergy, consumption.ExportedEnergy)

			aggregation[month] = aggData
		}
	}

	return aggregation
}
