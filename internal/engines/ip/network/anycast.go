package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type manycastResponse struct {
	Prefix     string `json:"prefix"`
	Anycast    bool   `json:"anycast"`
	Confidence string `json:"confidence"`

	ASNs []uint32 `json:"asns"`

	BackingPrefix string `json:"backing_prefix"`

	ABICMP int `json:"ab_icmp"`
	ABTCP  int `json:"ab_tcp"`
	ABDNS  int `json:"ab_dns"`

	GCDICMP int `json:"gcd_icmp"`
	GCDTCP  int `json:"gcd_tcp"`

	Locations []struct {
		City    *string `json:"city"`
		Country *string `json:"country"`

		ID string `json:"id"`

		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"locations"`

	Date string `json:"date"`

	QueriedIP    string `json:"queried_ip"`
	MappedPrefix string `json:"mapped_prefix"`
}

func manycast(
	ip string,
) (AnycastResult, error) {
	requestURL :=
		"https://manycast.net/api/v1/ip/" +
			url.PathEscape(ip)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		return AnycastResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return AnycastResult{},
			fmt.Errorf(
				"manycast returned HTTP %d",
				resp.StatusCode,
			)
	}

	var data manycastResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&data); err != nil {
		return AnycastResult{}, err
	}

	result := AnycastResult{
		Available: true,
		Anycast:   data.Anycast,

		Confidence: data.Confidence,

		Prefix:        data.Prefix,
		BackingPrefix: data.BackingPrefix,
		MappedPrefix:  data.MappedPrefix,

		ASNs: data.ASNs,

		ABICMP: data.ABICMP,
		ABTCP:  data.ABTCP,
		ABDNS:  data.ABDNS,

		GCDICMP: data.GCDICMP,
		GCDTCP:  data.GCDTCP,

		Date: data.Date,

		Source: requestURL,
	}

	for _, location := range data.Locations {
		city := ""
		country := ""

		if location.City != nil {
			city = *location.City
		}

		if location.Country != nil {
			country = *location.Country
		}

		result.Locations = append(
			result.Locations,
			AnycastLocation{
				ID: location.ID,

				City:    city,
				Country: country,

				Latitude:  location.Latitude,
				Longitude: location.Longitude,
			},
		)
	}

	return result, nil
}
