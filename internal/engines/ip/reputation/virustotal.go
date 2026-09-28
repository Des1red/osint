package reputation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type virusTotalResponse struct {
	Data struct {
		Attributes struct {
			Reputation int `json:"reputation"`

			LastAnalysisStats struct {
				Harmless   int `json:"harmless"`
				Malicious  int `json:"malicious"`
				Suspicious int `json:"suspicious"`
				Undetected int `json:"undetected"`
				Timeout    int `json:"timeout"`
			} `json:"last_analysis_stats"`

			TotalVotes struct {
				Harmless  int `json:"harmless"`
				Malicious int `json:"malicious"`
			} `json:"total_votes"`

			Tags []string `json:"tags"`

			LastAnalysisDate int64 `json:"last_analysis_date"`
		} `json:"attributes"`
	} `json:"data"`
}

func virusTotal(
	ip string,
) (VirusTotalResult, bool, error) {
	key :=
		os.Getenv(
			"VIRUSTOTAL_API_KEY",
		)

	if key == "" {
		return VirusTotalResult{},
			false,
			nil
	}

	requestURL :=
		"https://www.virustotal.com/api/v3/ip_addresses/" +
			url.PathEscape(ip)

	request, err :=
		http.NewRequest(
			http.MethodGet,
			requestURL,
			nil,
		)

	if err != nil {
		return VirusTotalResult{},
			true,
			err
	}

	request.Header.Set(
		"x-apikey",
		key,
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err :=
		client.Do(
			request,
		)

	if err != nil {
		return VirusTotalResult{},
			true,
			err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return VirusTotalResult{},
			true,
			fmt.Errorf(
				"HTTP %d",
				response.StatusCode,
			)
	}

	var data virusTotalResponse

	err =
		json.NewDecoder(
			response.Body,
		).Decode(
			&data,
		)

	if err != nil {
		return VirusTotalResult{},
			true,
			err
	}

	attributes :=
		data.Data.Attributes

	return VirusTotalResult{
			Available: true,

			Reputation: attributes.Reputation,

			Stats: VirusTotalStats{
				Harmless: attributes.LastAnalysisStats.Harmless,

				Malicious: attributes.LastAnalysisStats.Malicious,

				Suspicious: attributes.LastAnalysisStats.Suspicious,

				Undetected: attributes.LastAnalysisStats.Undetected,

				Timeout: attributes.LastAnalysisStats.Timeout,
			},

			Votes: VirusTotalVotes{
				Harmless: attributes.TotalVotes.Harmless,

				Malicious: attributes.TotalVotes.Malicious,
			},

			Tags: attributes.Tags,

			LastAnalysisDate: attributes.LastAnalysisDate,

			Source: requestURL,
		},
		true,
		nil
}
