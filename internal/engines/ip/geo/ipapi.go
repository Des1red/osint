package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ipAPIResponse struct {
	IP string `json:"ip"`

	CountryName string `json:"country_name"`
	CountryCode string `json:"country_code"`

	Region     string `json:"region"`
	RegionCode string `json:"region_code"`

	City   string `json:"city"`
	Postal string `json:"postal"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Timezone string `json:"timezone"`

	ContinentCode string `json:"continent_code"`

	Error  bool   `json:"error"`
	Reason string `json:"reason"`
}

func ipAPI(ip string) (ProviderResult, error) {
	url := fmt.Sprintf(
		"https://ipapi.co/%s/json/",
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
			"ipapi returned HTTP %d",
			resp.StatusCode,
		)
	}

	var data ipAPIResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ProviderResult{}, err
	}

	if data.Error {
		return ProviderResult{}, fmt.Errorf(
			"ipapi error: %s",
			data.Reason,
		)
	}

	return ProviderResult{
		Provider: "ipapi.co",

		IP: data.IP,

		Country:     data.CountryName,
		CountryCode: data.CountryCode,

		Region:     data.Region,
		RegionCode: data.RegionCode,

		City:       data.City,
		PostalCode: data.Postal,

		Latitude:  data.Latitude,
		Longitude: data.Longitude,

		Timezone: data.Timezone,

		ContinentCode: data.ContinentCode,

		Source: url,
	}, nil
}
