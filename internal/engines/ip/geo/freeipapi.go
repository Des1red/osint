package geo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type freeIPAPIResponse struct {
	IPVersion int    `json:"ipVersion"`
	IPAddress string `json:"ipAddress"`

	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	CountryName string `json:"countryName"`
	CountryCode string `json:"countryCode"`

	RegionName string `json:"regionName"`
	RegionCode string `json:"regionCode"`

	CityName string `json:"cityName"`
	ZipCode  string `json:"zipCode"`

	Continent     string `json:"continent"`
	ContinentCode string `json:"continentCode"`

	TimeZones []string `json:"timeZones"`
}

func freeIPAPI(ip string) (ProviderResult, error) {
	url := fmt.Sprintf(
		"https://free.freeipapi.com/api/v1/json/%s",
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
			"freeipapi returned HTTP %d",
			resp.StatusCode,
		)
	}

	var data freeIPAPIResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return ProviderResult{}, err
	}

	timezone := ""

	if len(data.TimeZones) == 1 {
		timezone = data.TimeZones[0]
	}

	return ProviderResult{
		Provider: "freeipapi.com",

		IP: data.IPAddress,

		Country:     data.CountryName,
		CountryCode: data.CountryCode,

		Region:     data.RegionName,
		RegionCode: data.RegionCode,

		City:       data.CityName,
		PostalCode: data.ZipCode,

		Latitude:  data.Latitude,
		Longitude: data.Longitude,

		Timezone: timezone,

		Continent:     data.Continent,
		ContinentCode: data.ContinentCode,

		Source: url,
	}, nil
}
