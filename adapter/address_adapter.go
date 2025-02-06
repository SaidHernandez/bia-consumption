package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type AddressAdapterInterface interface {
	GetAddress(meterID int) (*Address, error)
}

type Address struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

type AddressAdapter struct{}

var MockAddress = Address{
	ID:      0,
	Address: "Dirección Mock",
}

func NewAddressAdapter() *AddressAdapter {
	return &AddressAdapter{}
}

func callAddressService(meterID int) (Address, error) {
	url := fmt.Sprintf("http://localhost:8082/address/%d", meterID)

	resp, err := http.Get(url)
	if err != nil {
		return Address{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var address Address
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return Address{}, err
		}

		if err := json.Unmarshal(body, &address); err != nil {
			return Address{}, err
		}

		return address, nil
	}

	return Address{}, errors.New("failed to fetch address")
}

func (a *AddressAdapter) GetAddress(meterID int) (*Address, error) {
	var address Address
	var err error

	for attempts := 0; attempts < 2; attempts++ {
		address, err = callAddressService(meterID)
		if err == nil {
			return &address, nil
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return &MockAddress, nil
	}

	return &Address{}, err
}
