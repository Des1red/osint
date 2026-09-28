package network

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type ripeNetworkInfoResponse struct {
	Status string `json:"status"`

	Data struct {
		ASNs   []string `json:"asns"`
		Prefix string   `json:"prefix"`
	} `json:"data"`
}

type ripeASOverviewResponse struct {
	Status string `json:"status"`

	Data struct {
		Resource  string `json:"resource"`
		Holder    string `json:"holder"`
		Announced bool   `json:"announced"`
	} `json:"data"`
}

type ripeAnnouncedPrefixesResponse struct {
	Status string `json:"status"`

	Data struct {
		Prefixes []struct {
			Prefix string `json:"prefix"`
		} `json:"prefixes"`
	} `json:"data"`
}

type ripeASNNeighboursResponse struct {
	Status string `json:"status"`

	Data struct {
		NeighbourCounts struct {
			Left      int `json:"left"`
			Right     int `json:"right"`
			Uncertain int `json:"uncertain"`
			Unique    int `json:"unique"`
		} `json:"neighbour_counts"`

		Neighbours []struct {
			ASN uint32 `json:"asn"`

			Position string `json:"position"`
			Type     string `json:"type"`

			PathCount int `json:"path_count"`
			Power     int `json:"power"`

			PeerCount struct {
				V4 int `json:"v4"`
				V6 int `json:"v6"`
			} `json:"peer_count"`

			V4Peers int `json:"v4_peers"`
			V6Peers int `json:"v6_peers"`
		} `json:"neighbours"`
	} `json:"data"`
}

func ripeNetworkInfo(
	ip string,
) (string, []uint32, error) {
	endpoint :=
		"https://stat.ripe.net/data/network-info/data.json"

	requestURL := endpoint +
		"?resource=" +
		url.QueryEscape(ip)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf(
			"RIPEstat network-info returned HTTP %d",
			resp.StatusCode,
		)
	}

	var data ripeNetworkInfoResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&data); err != nil {
		return "", nil, err
	}

	if data.Status != "ok" {
		return "", nil, fmt.Errorf(
			"RIPEstat returned status %q",
			data.Status,
		)
	}

	var asns []uint32

	for _, value := range data.Data.ASNs {
		asn, err := strconv.ParseUint(
			value,
			10,
			32,
		)

		if err != nil {
			continue
		}

		asns = append(
			asns,
			uint32(asn),
		)
	}

	return strings.TrimSpace(
			data.Data.Prefix,
		),
		asns,
		nil
}

func ripeASOverview(
	asn uint32,
) (ASNResult, error) {
	endpoint :=
		"https://stat.ripe.net/data/as-overview/data.json"

	resource := asnResource(asn)

	requestURL := endpoint +
		"?resource=" +
		url.QueryEscape(resource)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		return ASNResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ASNResult{}, fmt.Errorf(
			"RIPEstat AS overview returned HTTP %d",
			resp.StatusCode,
		)
	}

	var data ripeASOverviewResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&data); err != nil {
		return ASNResult{}, err
	}

	if data.Status != "ok" {
		return ASNResult{}, fmt.Errorf(
			"RIPEstat returned status %q",
			data.Status,
		)
	}

	return ASNResult{
		ASN:       asn,
		Holder:    data.Data.Holder,
		Announced: data.Data.Announced,
	}, nil
}

func ripeAnnouncedPrefixes(
	asn uint32,
) ([]string, error) {
	endpoint :=
		"https://stat.ripe.net/data/announced-prefixes/data.json"

	requestURL := endpoint +
		"?resource=" +
		url.QueryEscape(
			asnResource(asn),
		)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"RIPEstat announced-prefixes returned HTTP %d",
			resp.StatusCode,
		)
	}

	var data ripeAnnouncedPrefixesResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&data); err != nil {
		return nil, err
	}

	if data.Status != "ok" {
		return nil, fmt.Errorf(
			"RIPEstat announced-prefixes returned status %q",
			data.Status,
		)
	}

	prefixes := make(
		[]string,
		0,
		len(data.Data.Prefixes),
	)

	for _, prefix := range data.Data.Prefixes {
		value := strings.TrimSpace(
			prefix.Prefix,
		)

		if value == "" {
			continue
		}

		prefixes = appendUnique(
			prefixes,
			value,
		)
	}

	sort.Strings(prefixes)

	return prefixes, nil
}

func ripeASNNeighbours(
	asn uint32,
) (
	ASNNeighbourCounts,
	[]ASNNeighbour,
	error,
) {
	endpoint :=
		"https://stat.ripe.net/data/asn-neighbours/data.json"

	requestURL := endpoint +
		"?resource=" +
		url.QueryEscape(
			asnResource(asn),
		) +
		"&lod=1"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(requestURL)
	if err != nil {
		return ASNNeighbourCounts{},
			nil,
			err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ASNNeighbourCounts{},
			nil,
			fmt.Errorf(
				"RIPEstat ASN neighbours returned HTTP %d",
				resp.StatusCode,
			)
	}

	var data ripeASNNeighboursResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&data); err != nil {
		return ASNNeighbourCounts{},
			nil,
			err
	}

	if data.Status != "ok" {
		return ASNNeighbourCounts{},
			nil,
			fmt.Errorf(
				"RIPEstat ASN neighbours returned status %q",
				data.Status,
			)
	}

	counts := ASNNeighbourCounts{
		Left: data.Data.NeighbourCounts.Left,

		Right: data.Data.NeighbourCounts.Right,

		Uncertain: data.Data.NeighbourCounts.Uncertain,

		Unique: data.Data.NeighbourCounts.Unique,
	}

	neighbours := make(
		[]ASNNeighbour,
		0,
		len(data.Data.Neighbours),
	)

	for _, neighbour := range data.Data.Neighbours {

		position :=
			strings.TrimSpace(
				neighbour.Position,
			)

		// RIPEstat older endpoint versions
		// use "type" instead of "position".
		if position == "" {
			position =
				strings.TrimSpace(
					neighbour.Type,
				)
		}

		pathCount :=
			neighbour.PathCount

		// Older RIPEstat responses expose
		// the same concept as "power".
		if pathCount == 0 {
			pathCount =
				neighbour.Power
		}

		peerCountV4 :=
			neighbour.PeerCount.V4

		peerCountV6 :=
			neighbour.PeerCount.V6

		if peerCountV4 == 0 {
			peerCountV4 =
				neighbour.V4Peers
		}

		if peerCountV6 == 0 {
			peerCountV6 =
				neighbour.V6Peers
		}

		neighbours = append(
			neighbours,
			ASNNeighbour{
				ASN: neighbour.ASN,

				Position: position,

				PathCount: pathCount,

				PeerCountV4: peerCountV4,

				PeerCountV6: peerCountV6,
			},
		)
	}

	sort.Slice(
		neighbours,
		func(i, j int) bool {
			if neighbours[i].Position !=
				neighbours[j].Position {

				return neighbours[i].Position <
					neighbours[j].Position
			}

			return neighbours[i].ASN <
				neighbours[j].ASN
		},
	)

	return counts,
		neighbours,
		nil
}

func asnResource(
	asn uint32,
) string {
	return "AS" +
		strconv.FormatUint(
			uint64(asn),
			10,
		)
}
