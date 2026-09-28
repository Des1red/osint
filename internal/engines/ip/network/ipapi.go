package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ipAPIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`

	Query string `json:"query"`

	ISP string `json:"isp"`
	Org string `json:"org"`

	AS     string `json:"as"`
	ASName string `json:"asname"`

	Reverse string `json:"reverse"`

	Mobile  bool `json:"mobile"`
	Proxy   bool `json:"proxy"`
	Hosting bool `json:"hosting"`
}

func ipAPI(
	ip string,
) (ipAPIResponse, error) {
	fields :=
		"status,message,query,isp,org,as,asname,reverse,mobile,proxy,hosting"

	requestURL :=
		"http://ip-api.com/json/" +
			url.PathEscape(ip) +
			"?fields=" +
			url.QueryEscape(fields)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		return ipAPIResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ipAPIResponse{},
			fmt.Errorf(
				"ip-api.com returned HTTP %d",
				resp.StatusCode,
			)
	}

	var data ipAPIResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&data); err != nil {
		return ipAPIResponse{}, err
	}

	if data.Status != "success" {
		return ipAPIResponse{},
			fmt.Errorf(
				"ip-api.com: %s",
				data.Message,
			)
	}

	return data, nil
}
