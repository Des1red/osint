package history

import (
	"net/url"
	"sort"
)

type routingHistoryResponse struct {
	Resource string `json:"resource"`

	QueryStartTime string `json:"query_starttime"`

	QueryEndTime string `json:"query_endtime"`

	ByOrigin []struct {
		Origin string `json:"origin"`

		Prefixes []struct {
			Prefix string `json:"prefix"`

			Timelines []struct {
				StartTime string `json:"starttime"`

				EndTime string `json:"endtime"`

				FullPeersSeeing int `json:"full_peers_seeing"`

				Visibility *float64 `json:"visibility"`
			} `json:"timelines"`
		} `json:"prefixes"`
	} `json:"by_origin"`
}

func routingHistory(
	ip string,
) (
	RoutingHistoryResult,
	error,
) {
	values :=
		url.Values{}

	values.Set(
		"resource",
		ip,
	)

	//
	// Keep the response bounded while still
	// collecting meaningful routing history.
	//
	values.Set(
		"max_rows",
		"500",
	)

	//
	// Include normalized RIS visibility for
	// each historical routing period.
	//
	values.Set(
		"normalise_visibility",
		"true",
	)

	var data routingHistoryResponse

	source,
		err :=
		ripeStatGet(
			"routing-history",
			values,
			&data,
		)

	if err != nil {
		return RoutingHistoryResult{},
			err
	}

	result :=
		RoutingHistoryResult{
			Available: true,

			Resource: data.Resource,

			QueryStartTime: data.QueryStartTime,

			QueryEndTime: data.QueryEndTime,

			Source: source,
		}

	for _, origin := range data.ByOrigin {

		originResult :=
			RoutingOrigin{
				Origin: origin.Origin,
			}

		for _, prefix := range origin.Prefixes {

			prefixResult :=
				RoutingPrefix{
					Prefix: prefix.Prefix,
				}

			for _, timeline := range prefix.Timelines {

				timelineResult :=
					RoutingTimeline{
						StartTime: timeline.StartTime,

						EndTime: timeline.EndTime,

						FullPeersSeeing: timeline.FullPeersSeeing,
					}

				if timeline.Visibility !=
					nil {

					timelineResult.Visibility =
						*timeline.Visibility

					timelineResult.VisibilityAvailable =
						true
				}

				prefixResult.Timelines =
					append(
						prefixResult.Timelines,
						timelineResult,
					)
			}

			//
			// Most recent historical periods
			// first.
			//
			sort.SliceStable(
				prefixResult.Timelines,

				func(
					i int,
					j int,
				) bool {
					return prefixResult.
						Timelines[i].
						StartTime >
						prefixResult.
							Timelines[j].
							StartTime
				},
			)

			originResult.Prefixes =
				append(
					originResult.Prefixes,
					prefixResult,
				)
		}

		result.Origins =
			append(
				result.Origins,
				originResult,
			)
	}

	return result,
		nil
}
