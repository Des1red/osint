package history

import (
	"encoding/json"
	"net/url"
	"strings"
)

type routingSeenResponse struct {
	Time string `json:"time"`

	Origin json.RawMessage `json:"origin"`

	Prefix string `json:"prefix"`
}

type routingStatusResponse struct {
	Resource string `json:"resource"`

	QueryTime string `json:"query_time"`

	FirstSeen routingSeenResponse `json:"first_seen"`

	LastSeen routingSeenResponse `json:"last_seen"`
}

func routingStatus(
	ip string,
) (
	RoutingStatusResult,
	error,
) {
	values :=
		url.Values{}

	values.Set(
		"resource",
		ip,
	)

	var data routingStatusResponse

	source,
		err :=
		ripeStatGet(
			"routing-status",
			values,
			&data,
		)

	if err != nil {
		return RoutingStatusResult{},
			err
	}

	result :=
		RoutingStatusResult{
			Available: true,

			Resource: data.Resource,

			QueryTime: data.QueryTime,

			FirstSeen: seenResult(
				data.FirstSeen,
			),

			LastSeen: seenResult(
				data.LastSeen,
			),

			Source: source,
		}

	return result,
		nil
}

func seenResult(
	input routingSeenResponse,
) SeenResult {
	origin :=
		rawOrigin(
			input.Origin,
		)

	available :=
		strings.TrimSpace(
			input.Time,
		) != "" ||
			strings.TrimSpace(
				input.Prefix,
			) != "" ||
			origin != ""

	return SeenResult{
		Available: available,

		Time: input.Time,

		Origin: origin,

		Prefix: input.Prefix,
	}
}

func rawOrigin(
	raw json.RawMessage,
) string {
	if len(raw) == 0 {
		return ""
	}

	var stringValue string

	if err :=
		json.Unmarshal(
			raw,
			&stringValue,
		); err == nil {

		return strings.TrimSpace(
			stringValue,
		)
	}

	var numberValue json.Number

	if err :=
		json.Unmarshal(
			raw,
			&numberValue,
		); err == nil {

		return numberValue.String()
	}

	return ""
}
