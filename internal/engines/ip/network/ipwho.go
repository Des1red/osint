package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ipWhoResponse struct {
	Success bool `json:"success"`

	Data struct {
		IP string `json:"ip"`

		Connection struct {
			ASNNumber uint32 `json:"asn_number"`
			ASNOrg    string `json:"asn_org"`

			ISP string `json:"isp"`
			Org string `json:"org"`

			Domain string `json:"domain"`

			ConnectionType string `json:"connection_type"`
		} `json:"connection"`

		Security struct {
			IsVPN    bool   `json:"isVpn"`
			IsTor    bool   `json:"isTor"`
			IsThreat string `json:"isThreat"`

			IsVPNSnake    bool   `json:"is_vpn"`
			IsTorSnake    bool   `json:"is_tor"`
			IsThreatSnake string `json:"is_threat"`
		} `json:"security"`
	} `json:"data"`
}

type ipWhoResult struct {
	ASN            uint32
	ASNOrg         string
	ISP            string
	Organization   string
	Domain         string
	ConnectionType string

	VPN    bool
	Tor    bool
	Threat string
}

func ipWho(
	ip string,
) (ipWhoResult, error) {
	key := os.Getenv(
		"IPWHO_API_KEY",
	)

	if key == "" {
		return ipWhoResult{},
			fmt.Errorf(
				"IPWHO_API_KEY not configured",
			)
	}

	requestURL :=
		"https://api.ipwho.org/ip/" +
			url.PathEscape(ip) +
			"?apiKey=" +
			url.QueryEscape(key)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		return ipWhoResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ipWhoResult{}, fmt.Errorf(
			"IPWho returned HTTP %d",
			resp.StatusCode,
		)
	}

	var data ipWhoResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&data); err != nil {
		return ipWhoResult{}, err
	}

	if !data.Success {
		return ipWhoResult{},
			fmt.Errorf(
				"IPWho lookup failed",
			)
	}

	vpn := data.Data.Security.IsVPN ||
		data.Data.Security.IsVPNSnake

	tor := data.Data.Security.IsTor ||
		data.Data.Security.IsTorSnake

	threat := data.Data.Security.IsThreat

	if threat == "" {
		threat =
			data.Data.Security.IsThreatSnake
	}

	return ipWhoResult{
		ASN: data.Data.Connection.ASNNumber,

		ASNOrg: data.Data.Connection.ASNOrg,

		ISP: data.Data.Connection.ISP,

		Organization: data.Data.Connection.Org,

		Domain: data.Data.Connection.Domain,

		ConnectionType: data.Data.Connection.ConnectionType,

		VPN: vpn,
		Tor: tor,

		Threat: threat,
	}, nil
}
