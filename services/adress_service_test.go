package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/SaidHernandez/bia-comsumtion/adapter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMemoryCache struct {
	mock.Mock
}

func (m *MockMemoryCache) Get(ctx context.Context, key string) (interface{}, bool, error) {
	args := m.Called(ctx, key)
	return args.Get(0), args.Bool(1), args.Error(2)
}

func (m *MockMemoryCache) Set(ctx context.Context, key string, value interface{}, expireAfter time.Duration) error {
	args := m.Called(ctx, key, value, expireAfter)
	return args.Error(0)
}

func (m *MockMemoryCache) Clear(ctx context.Context, key string) (bool, error) {
	args := m.Called(ctx, key)
	return args.Bool(0), args.Error(1)
}

type MockAdapter struct {
	mock.Mock
}

func (m *MockAdapter) GetAddress(meterId int) (*adapter.Address, error) {
	args := m.Called(meterId)
	return args.Get(0).(*adapter.Address), args.Error(1)
}

func TestAddressService_GetAddress(t *testing.T) {
	tests := []struct {
		name          string
		meterId       int
		mockCache     func() *MockMemoryCache
		mockAdapter   func() *MockAdapter
		expectedError error
		expectedAddr  *adapter.Address
	}{
		{
			name:    "Success: Get address from cache",
			meterId: 1,
			mockCache: func() *MockMemoryCache {
				cacheMock := new(MockMemoryCache)
				cacheMock.On("Get", mock.Anything, "address-1").Return(&adapter.Address{ID: 1, Address: "Main St"}, true, nil)
				return cacheMock
			},
			mockAdapter:   func() *MockAdapter { return nil },
			expectedError: nil,
			expectedAddr:  &adapter.Address{ID: 1, Address: "Main St"},
		},
		{
			name:    "Success: Get address from adapter and store in cache",
			meterId: 2,
			mockCache: func() *MockMemoryCache {
				cacheMock := new(MockMemoryCache)
				cacheMock.On("Get", mock.Anything, "address-2").Return(nil, false, nil)
				cacheMock.On("Set", mock.Anything, "address-2", mock.Anything, time.Duration(24)*time.Hour).Return(nil)
				return cacheMock
			},
			mockAdapter: func() *MockAdapter {
				adapterMock := new(MockAdapter)
				adapterMock.On("GetAddress", 2).Return(&adapter.Address{ID: 2, Address: "Main St"}, nil)
				return adapterMock
			},
			expectedError: nil,
			expectedAddr:  &adapter.Address{ID: 2, Address: "Main St"},
		},
		{
			name:    "Error: Get address from cache fails",
			meterId: 3,
			mockCache: func() *MockMemoryCache {
				cacheMock := new(MockMemoryCache)
				cacheMock.On("Get", mock.Anything, "address-3").Return(nil, false, errors.New("cache error"))
				return cacheMock
			},
			mockAdapter:   func() *MockAdapter { return nil },
			expectedError: errors.New("error al obtener la dirección desde la caché: cache error"),
			expectedAddr:  nil,
		},
		{
			name:    "Error: Get address from adapter fails",
			meterId: 4,
			mockCache: func() *MockMemoryCache {
				cacheMock := new(MockMemoryCache)
				cacheMock.On("Get", mock.Anything, "address-4").Return(nil, false, nil)
				return cacheMock
			},
			mockAdapter: func() *MockAdapter {
				adapterMock := new(MockAdapter)
				// Devolver un puntero nulo explícito
				adapterMock.On("GetAddress", 4).Return((*adapter.Address)(nil), errors.New("adapter error"))
				return adapterMock
			},
			expectedError: errors.New("no se pudo obtener la dirección después de varios intentos: adapter error"),
			expectedAddr:  nil,
		},
		{
			name:    "Error: Cache store fails",
			meterId: 5,
			mockCache: func() *MockMemoryCache {
				cacheMock := new(MockMemoryCache)
				cacheMock.On("Get", mock.Anything, "address-5").Return(nil, false, nil)
				cacheMock.On("Set", mock.Anything, "address-5", mock.Anything, time.Duration(24)*time.Hour).Return(errors.New("cache store error"))
				return cacheMock
			},
			mockAdapter: func() *MockAdapter {
				adapterMock := new(MockAdapter)
				adapterMock.On("GetAddress", 5).Return(&adapter.Address{ID: 2, Address: "Main St"}, nil)
				return adapterMock
			},
			expectedError: errors.New("error al almacenar la dirección en la caché: cache store error"),
			expectedAddr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cacheMock := tt.mockCache()
			adapterMock := tt.mockAdapter()
			service := NewAddressServiceClient(cacheMock, adapterMock)

			address, err := service.GetAddress(context.Background(), tt.meterId)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedAddr, address)
			}

			cacheMock.AssertExpectations(t)
			if adapterMock != nil {
				adapterMock.AssertExpectations(t)
			}
		})
	}
}
