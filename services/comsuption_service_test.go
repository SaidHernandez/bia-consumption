package services

import (
	"context"
	"testing"
	"time"

	"github.com/SaidHernandez/bia-comsumtion/adapter"
	"github.com/SaidHernandez/bia-comsumtion/business/model"
	"github.com/SaidHernandez/bia-comsumtion/business/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAddressService struct {
	mock.Mock
}

func (m *MockAddressService) GetAddress(ctx context.Context, meterID int) (*adapter.Address, error) {
	args := m.Called(ctx, meterID)
	return args.Get(0).(*adapter.Address), args.Error(1)
}

var _ AddressServiceInterface = (*MockAddressService)(nil)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetConsumptionByFilters(meterID int, startDate, endDate string) ([]model.Consumption, error) {
	args := m.Called(meterID, startDate, endDate)
	return args.Get(0).([]model.Consumption), args.Error(1)
}

var _ repository.ConsumptionRepositoryInterface = (*MockRepository)(nil)

func TestConsumptionService_GetConsumptionByPeriod(t *testing.T) {
	tests := []struct {
		name            string
		meterIDs        []int
		startDate       string
		endDate         string
		kindPeriod      string
		mockAddress     func() AddressServiceInterface
		mockRepository  func() repository.ConsumptionRepositoryInterface
		expectedResults map[string]interface{}
		expectedError   error
	}{
		{
			name:       "Success: Get monthly consumption data",
			meterIDs:   []int{1, 2},
			startDate:  "2023-06-01",
			endDate:    "2023-06-30",
			kindPeriod: "monthly",
			mockAddress: func() AddressServiceInterface {
				addressMock := new(MockAddressService)
				addressMock.On("GetAddress", mock.Anything, 1).Return(&adapter.Address{Address: "123 Main St"}, nil)
				addressMock.On("GetAddress", mock.Anything, 2).Return(&adapter.Address{Address: "456 Side St"}, nil)
				return addressMock
			},
			mockRepository: func() repository.ConsumptionRepositoryInterface {
				repoMock := new(MockRepository)
				dateStr := "2023-07-04 10:59:00+00"
				date, _ := time.Parse("2006-01-02 15:04:05-07", dateStr)
				repoMock.On("GetConsumptionByFilters", 1, "2023-06-01", "2023-06-30").Return([]model.Consumption{
					{ID: "1", MeterID: 1, ActiveEnergy: 100, ReactiveInductive: 50, ReactiveCapacitive: 0, ExportedEnergy: 0, Date: date},
				}, nil)

				repoMock.On("GetConsumptionByFilters", 2, "2023-06-01", "2023-06-30").Return([]model.Consumption{
					{ID: "2", MeterID: 2, ActiveEnergy: 200, ReactiveInductive: 100, ReactiveCapacitive: 0, ExportedEnergy: 0, Date: date},
				}, nil)

				return repoMock
			},
			expectedResults: map[string]interface{}{
				"period": []string{"Jul 2023"},
				"data_graph": []map[string]interface{}{
					{
						"active":              []float64{200},
						"reactive_inductive":  []float64{100},
						"reactive_capacitive": []float64{0},
						"exported":            []float64{0},
						"address":             "456 Side St",
						"meter_id":            2,
					},
					{
						"active":              []float64{100},
						"reactive_inductive":  []float64{50},
						"reactive_capacitive": []float64{0},
						"exported":            []float64{0},
						"address":             "123 Main St",
						"meter_id":            1,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:       "Success: Get weekly consumption data",
			meterIDs:   []int{1, 2},
			startDate:  "2023-06-01",
			endDate:    "2023-06-30",
			kindPeriod: "weekly",
			mockAddress: func() AddressServiceInterface {
				addressMock := new(MockAddressService)
				addressMock.On("GetAddress", mock.Anything, 1).Return(&adapter.Address{Address: "123 Main St"}, nil)
				addressMock.On("GetAddress", mock.Anything, 2).Return(&adapter.Address{Address: "456 Side St"}, nil)
				return addressMock
			},
			mockRepository: func() repository.ConsumptionRepositoryInterface {
				repoMock := new(MockRepository)
				date1Str := "2023-06-03 10:59:00+00"
				date1, _ := time.Parse("2006-01-02 15:04:05-07", date1Str)
				date2Str := "2023-06-10 10:59:00+00"
				date2, _ := time.Parse("2006-01-02 15:04:05-07", date2Str)

				repoMock.On("GetConsumptionByFilters", 1, "2023-06-01", "2023-06-30").Return([]model.Consumption{
					{ID: "1", MeterID: 1, ActiveEnergy: 100, ReactiveInductive: 50, ReactiveCapacitive: 0, ExportedEnergy: 0, Date: date1},
					{ID: "1", MeterID: 1, ActiveEnergy: 150, ReactiveInductive: 70, ReactiveCapacitive: 0, ExportedEnergy: 0, Date: date2},
				}, nil)

				repoMock.On("GetConsumptionByFilters", 2, "2023-06-01", "2023-06-30").Return([]model.Consumption{
					{ID: "2", MeterID: 2, ActiveEnergy: 200, ReactiveInductive: 100, ReactiveCapacitive: 0, ExportedEnergy: 0, Date: date1},
					{ID: "2", MeterID: 2, ActiveEnergy: 250, ReactiveInductive: 120, ReactiveCapacitive: 0, ExportedEnergy: 0, Date: date2},
				}, nil)

				return repoMock
			},
			expectedResults: map[string]interface{}{
				"period": []string{
					"May 28 - Jun 3",
					"Jun 4 - Jun 10",
				},
				"data_graph": []map[string]interface{}{
					{
						"active":              []float64{200, 250},
						"reactive_inductive":  []float64{100, 120},
						"reactive_capacitive": []float64{0, 0},
						"exported":            []float64{0, 0},
						"address":             "456 Side St",
						"meter_id":            2,
					},
					{
						"active":              []float64{100, 150},
						"reactive_inductive":  []float64{50, 70},
						"reactive_capacitive": []float64{0, 0},
						"exported":            []float64{0, 0},
						"address":             "123 Main St",
						"meter_id":            1,
					},
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addressService := tt.mockAddress()
			repo := tt.mockRepository()
			service := NewConsumptionService(addressService, repo)

			results, err := service.GetConsumptionByPeriod(context.Background(), tt.meterIDs, tt.startDate, tt.endDate, tt.kindPeriod)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResults, results)
			}
		})
	}
}
