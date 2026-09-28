package history

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const ripeStatBaseURL = "https://stat.ripe.net/data/"

const ripeStatSourceApp = "osint-master"

var ripeStatClient = &http.Client{
	Timeout: 15 *
		time.Second,
}

type ripeStatEnvelope struct {
	Status string `json:"status"`

	Message string `json:"message"`

	Data json.RawMessage `json:"data"`
}

func ripeStatGet(
	endpoint string,
	values url.Values,
	target any,
) (
	string,
	error,
) {
	values.Set(
		"sourceapp",
		ripeStatSourceApp,
	)

	requestURL :=
		ripeStatBaseURL +
			endpoint +
			"/data.json?" +
			values.Encode()

	request,
		err :=
		http.NewRequest(
			http.MethodGet,
			requestURL,
			nil,
		)

	if err != nil {
		return "",
			err
	}

	request.Header.Set(
		"Accept",
		"application/json",
	)

	request.Header.Set(
		"User-Agent",
		"OSINT-Master/1.0",
	)

	response,
		err :=
		ripeStatClient.Do(
			request,
		)

	if err != nil {
		return "",
			err
	}

	defer response.Body.Close()

	if response.StatusCode !=
		http.StatusOK {

		return "",
			fmt.Errorf(
				"RIPEstat %s returned HTTP %d",
				endpoint,
				response.StatusCode,
			)
	}

	var envelope ripeStatEnvelope

	err =
		json.NewDecoder(
			response.Body,
		).Decode(
			&envelope,
		)

	if err != nil {
		return "",
			err
	}

	if envelope.Status != "ok" {

		message :=
			envelope.Message

		if message == "" {
			message =
				"unknown RIPEstat error"
		}

		return "",
			fmt.Errorf(
				"RIPEstat %s: %s",
				endpoint,
				message,
			)
	}

	err =
		json.Unmarshal(
			envelope.Data,
			target,
		)

	if err != nil {
		return "",
			err
	}

	return requestURL,
		nil
}
