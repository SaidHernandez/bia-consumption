package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/SaidHernandez/bia-comsumtion/business/aggregate"
	"github.com/SaidHernandez/bia-comsumtion/business/repository"
)

type ConsumptionService struct {
	addressService AddressServiceInterface
	repository     repository.ConsumptionRepositoryInterface
}

func NewConsumptionService(addressService AddressServiceInterface, repository repository.ConsumptionRepositoryInterface) *ConsumptionService {
	return &ConsumptionService{
		addressService: addressService,
		repository:     repository,
	}
}

func (service *ConsumptionService) GetConsumptionByPeriod(ctx context.Context, meterIDs []int, startDate, endDate, kindPeriod string) (map[string]interface{}, error) {
	strategies := map[string]aggregate.AggregationStrategy{
		"monthly": &aggregate.MonthlyAggregationStrategy{},
		"weekly":  &aggregate.WeeklyAggregationStrategy{},
		"daily":   &aggregate.DailyAggregationStrategy{},
	}

	strategy, exists := strategies[kindPeriod]
	if !exists {
		return nil, fmt.Errorf("invalid kind_period: %s", kindPeriod)
	}

	var wg sync.WaitGroup
	var dataGraph []map[string]interface{}
	var periods []string
	mu := sync.Mutex{}

	for _, meterID := range meterIDs {
		wg.Add(1)
		go func(meterID int) {
			defer wg.Done()

			consumptions, err := service.repository.GetConsumptionByFilters(meterID, startDate, endDate)
			if err != nil {
				fmt.Println("Error fetching data for meterID", meterID, ":", err)
				return
			}

			aggregatedData := strategy.Aggregate(consumptions)

			address, err := service.addressService.GetAddress(ctx, meterID)
			if err != nil {
				fmt.Println("Error fetching address for meterID", meterID, ":", err)
				return
			}

			var active []float64
			var reactiveInductive []float64
			var reactiveCapacitive []float64
			var exported []float64
			localPeriods := []string{}

			for _, aggData := range aggregatedData {
				localPeriods = append(localPeriods, aggData.Period...)
				active = append(active, aggData.ActiveEnergy...)
				reactiveInductive = append(reactiveInductive, aggData.ReactiveInductive...)
				reactiveCapacitive = append(reactiveCapacitive, aggData.ReactiveCapacitive...)
				exported = append(exported, aggData.ExportedEnergy...)
			}

			mu.Lock()
			if len(periods) == 0 { // Guardar los periodos solo una vez
				periods = localPeriods
			}
			dataGraph = append(dataGraph, map[string]interface{}{
				"meter_id":            meterID,
				"address":             address.Address,
				"active":              active,
				"reactive_inductive":  reactiveInductive,
				"reactive_capacitive": reactiveCapacitive,
				"exported":            exported,
			})
			mu.Unlock()
		}(meterID)
	}

	wg.Wait()

	return map[string]interface{}{
		"period":     periods,
		"data_graph": dataGraph,
	}, nil
}
