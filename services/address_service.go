package services

import (
	"context"
	"fmt"
	"time"

	"github.com/SaidHernandez/bia-comsumtion/adapter"
	"github.com/SaidHernandez/bia-comsumtion/infraestructure/cache"
)

type AddressServiceInterface interface {
	GetAddress(ctx context.Context, meterID int) (*adapter.Address, error)
}

type AddressService struct {
	cache   cache.Cache
	adapter adapter.AddressAdapterInterface
}

func NewAddressServiceClient(cache cache.Cache, adapter adapter.AddressAdapterInterface) *AddressService {
	return &AddressService{
		cache:   cache,
		adapter: adapter,
	}
}

func (client *AddressService) GetAddress(ctx context.Context, meterId int) (*adapter.Address, error) {
	cacheKey := fmt.Sprintf("address-%d", meterId)

	address, found, err := client.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la dirección desde la caché: %w", err)
	}

	if found {
		return address.(*adapter.Address), nil
	}

	return client.getAddressAdapter(ctx, meterId)
}

func (client *AddressService) getAddressAdapter(ctx context.Context, meterId int) (*adapter.Address, error) {
	var address *adapter.Address
	var err error

	address, err = client.adapter.GetAddress(meterId)
	if err == nil {
		cacheKey := fmt.Sprintf("address-%d", meterId)
		err := client.cache.Set(ctx, cacheKey, address, 24*time.Hour)
		if err != nil {
			return nil, fmt.Errorf("error al almacenar la dirección en la caché: %w", err)
		}
		return address, nil
	}

	return nil, fmt.Errorf("no se pudo obtener la dirección después de varios intentos: %w", err)
}
