package aggregate

import (
	"fmt"
	"time"

	"github.com/SaidHernandez/bia-comsumtion/business/model"
)

type WeeklyAggregationStrategy struct{}

func (w *WeeklyAggregationStrategy) Aggregate(consumptions []model.Consumption) map[string]model.AggregatedConsumption {
	aggregation := make(map[string]model.AggregatedConsumption)

	formatWeekRange := func(date time.Time) string {
		startOfWeek := date.AddDate(0, 0, -int(date.Weekday()))
		endOfWeek := startOfWeek.AddDate(0, 0, 6)
		return fmt.Sprintf("%s %d - %s %d", startOfWeek.Format("Jan"), startOfWeek.Day(), endOfWeek.Format("Jan"), endOfWeek.Day())
	}

	for _, consumption := range consumptions {
		period := formatWeekRange(consumption.Date)

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
