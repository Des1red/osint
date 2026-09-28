package reputation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ipQSResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`

	FraudScore int `json:"fraud_score"`

	Proxy bool `json:"proxy"`
	VPN   bool `json:"vpn"`
	Tor   bool `json:"tor"`

	BotStatus      bool `json:"bot_status"`
	RecentAbuse    bool `json:"recent_abuse"`
	FrequentAbuser bool `json:"frequent_abuser"`

	AbuseVelocity string `json:"abuse_velocity"`
}

func ipQS(
	ip string,
) (IPQSResult, bool, error) {
	key :=
		os.Getenv(
			"IPQS_API_KEY",
		)

	if key == "" {
		return IPQSResult{},
			false,
			nil
	}

	requestURL :=
		"https://ipqualityscore.com/api/json/ip/" +
			url.PathEscape(key) +
			"/" +
			url.PathEscape(ip) +
			"?strictness=1"

	request, err :=
		http.NewRequest(
			http.MethodGet,
			requestURL,
			nil,
		)

	if err != nil {
		return IPQSResult{},
			true,
			err
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err :=
		client.Do(
			request,
		)

	if err != nil {
		return IPQSResult{},
			true,
			err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return IPQSResult{},
			true,
			fmt.Errorf(
				"HTTP %d",
				response.StatusCode,
			)
	}

	var data ipQSResponse

	err =
		json.NewDecoder(
			response.Body,
		).Decode(
			&data,
		)

	if err != nil {
		return IPQSResult{},
			true,
			err
	}

	if !data.Success {
		return IPQSResult{},
			true,
			fmt.Errorf(
				"%s",
				data.Message,
			)
	}

	return IPQSResult{
			Available: true,

			FraudScore: data.FraudScore,

			Proxy: data.Proxy,

			VPN: data.VPN,

			Tor: data.Tor,

			BotStatus: data.BotStatus,

			RecentAbuse: data.RecentAbuse,

			FrequentAbuser: data.FrequentAbuser,

			AbuseVelocity: data.AbuseVelocity,

			Source: "https://ipqualityscore.com/api/json/ip",
		},
		true,
		nil
}
