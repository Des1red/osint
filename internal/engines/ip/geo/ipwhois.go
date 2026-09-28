package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ipWhoIsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`

	IP string `json:"ip"`

	Continent     string `json:"continent"`
	ContinentCode string `json:"continent_code"`

	Country     string `json:"country"`
	CountryCode string `json:"country_code"`

	Region     string `json:"region"`
	RegionCode string `json:"region_code"`

	City   string `json:"city"`
	Postal string `json:"postal"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Timezone struct {
		ID string `json:"id"`
	} `json:"timezone"`
}

func ipWhoIs(ip string) (ProviderResult, error) {
	url := fmt.Sprintf(
		"https://ipwho.is/%s",
		ip,
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return ProviderResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ProviderResult{}, fmt.Errorf(
			"ipwho.is returned HTTP %d",
			resp.StatusCode,
		)
	}

	var data ipWhoIsResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ProviderResult{}, err
	}

	if !data.Success {
		return ProviderResult{}, fmt.Errorf(
			"ipwho.is error: %s",
			data.Message,
		)
	}

	return ProviderResult{
		Provider: "ipwho.is",

		IP: data.IP,

		Country:     data.Country,
		CountryCode: data.CountryCode,

		Region:     data.Region,
		RegionCode: data.RegionCode,

		City:       data.City,
		PostalCode: data.Postal,

		Latitude:  data.Latitude,
		Longitude: data.Longitude,

		Timezone: data.Timezone.ID,

		Continent:     data.Continent,
		ContinentCode: data.ContinentCode,

		Source: url,
	}, nil
}
