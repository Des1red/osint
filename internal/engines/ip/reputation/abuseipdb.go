package reputation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type abuseIPDBResponse struct {
	Data struct {
		IPAddress            string `json:"ipAddress"`
		IsPublic             bool   `json:"isPublic"`
		AbuseConfidenceScore int    `json:"abuseConfidenceScore"`
		UsageType            string `json:"usageType"`
		Domain               string `json:"domain"`
		TotalReports         int    `json:"totalReports"`
		NumDistinctUsers     int    `json:"numDistinctUsers"`
		LastReportedAt       string `json:"lastReportedAt"`

		IsWhitelisted *bool `json:"isWhitelisted"`
	} `json:"data"`
}

func abuseIPDB(
	ip string,
) (AbuseIPDBResult, bool, error) {
	key := os.Getenv(
		"ABUSEIPDB_API_KEY",
	)

	if key == "" {
		return AbuseIPDBResult{},
			false,
			nil
	}

	requestURL :=
		"https://api.abuseipdb.com/api/v2/check"

	values := url.Values{}

	values.Set(
		"ipAddress",
		ip,
	)

	values.Set(
		"maxAgeInDays",
		"90",
	)

	requestURL +=
		"?" +
			values.Encode()

	request, err :=
		http.NewRequest(
			http.MethodGet,
			requestURL,
			nil,
		)

	if err != nil {
		return AbuseIPDBResult{},
			true,
			err
	}

	request.Header.Set(
		"Key",
		key,
	)

	request.Header.Set(
		"Accept",
		"application/json",
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err :=
		client.Do(
			request,
		)

	if err != nil {
		return AbuseIPDBResult{},
			true,
			err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return AbuseIPDBResult{},
			true,
			fmt.Errorf(
				"HTTP %d",
				response.StatusCode,
			)
	}

	var data abuseIPDBResponse

	err =
		json.NewDecoder(
			response.Body,
		).Decode(
			&data,
		)

	if err != nil {
		return AbuseIPDBResult{},
			true,
			err
	}

	return AbuseIPDBResult{
			Available: true,

			AbuseConfidenceScore: data.Data.AbuseConfidenceScore,

			TotalReports: data.Data.TotalReports,

			DistinctUsers: data.Data.NumDistinctUsers,

			LastReportedAt: data.Data.LastReportedAt,

			UsageType: data.Data.UsageType,

			Domain: data.Data.Domain,

			IsPublic: data.Data.IsPublic,

			IsWhitelisted: data.Data.IsWhitelisted,

			Source: requestURL,
		},
		true,
		nil
}
