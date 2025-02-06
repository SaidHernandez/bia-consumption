package aggregate

import (
	"fmt"
	"time"

	"github.com/SaidHernandez/bia-comsumtion/business/model"
)

type DailyAggregationStrategy struct{}

func (d *DailyAggregationStrategy) Aggregate(consumptions []model.Consumption) map[string]model.AggregatedConsumption {
	aggregation := make(map[string]model.AggregatedConsumption)

	formatDay := func(date time.Time) string {
		return fmt.Sprintf("%s %d", date.Format("Jan"), date.Day())
	}

	for _, consumption := range consumptions {
		period := formatDay(consumption.Date)

		if _, exists := aggregation[period]; !exists {
			aggregation[period] = model.AggregatedConsumption{
				Period:             []string{period},
				ActiveEnergy:       []float64{},
				ReactiveInductive:  []float64{},
				ReactiveCapacitive: []float64{},
				ExportedEnergy:     []float64{},
			}
		}

		aggData := aggregation[period]

		aggData.ActiveEnergy = append(aggData.ActiveEnergy, consumption.ActiveEnergy)
		aggData.ReactiveInductive = append(aggData.ReactiveInductive, consumption.ReactiveInductive)
		aggData.ReactiveCapacitive = append(aggData.ReactiveCapacitive, consumption.ReactiveCapacitive)
		aggData.ExportedEnergy = append(aggData.ExportedEnergy, consumption.ExportedEnergy)

		aggregation[period] = aggData
	}

	return aggregation
}
