package subdomains

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"osint/internal/httpx"
)

const maxCandidates = 200

const maxResponseBytes int64 = 16 * 1024 * 1024

type Result struct {
	Source string

	Total int

	Names []string
}

type crtEntry struct {
	NameValue string `json:"name_value"`
}

func Discover(
	domain string,
) (
	Result,
	error,
) {
	requestURL :=
		crtURL(
			domain,
		)

	headers :=
		make(
			http.Header,
		)

	headers.Set(
		"Accept",
		"application/json",
	)

	headers.Set(
		"User-Agent",
		"osint-master",
	)

	response,
		err :=
		httpx.Do(
			httpx.Request{
				Method: http.MethodGet,

				URL: requestURL,

				Headers: headers,

				MaxBodyBytes: maxResponseBytes,
			},
		)

	if err != nil {

		return Result{
				Source: requestURL,
			},
			err
	}

	if response.StatusCode !=
		http.StatusOK {

		return Result{
				Source: requestURL,
			},
			fmt.Errorf(
				"crt.sh HTTP %d",
				response.StatusCode,
			)
	}

	if response.Truncated {

		return Result{
				Source: requestURL,
			},
			fmt.Errorf(
				"crt.sh response exceeded maximum body size",
			)
	}

	var entries []crtEntry

	err =
		json.Unmarshal(
			response.Body,
			&entries,
		)

	if err != nil {

		return Result{
				Source: requestURL,
			},
			err
	}

	names :=
		extractNames(
			domain,
			entries,
		)

	total :=
		len(
			names,
		)

	if len(names) >
		maxCandidates {

		names =
			names[:maxCandidates]
	}

	return Result{
			Source: requestURL,

			Total: total,

			Names: names,
		},
		nil
}

func crtURL(
	domain string,
) string {
	values :=
		url.Values{}

	values.Set(
		"q",
		"%."+domain,
	)

	values.Set(
		"output",
		"json",
	)

	return "https://crt.sh/?" +
		values.Encode()
}

func extractNames(
	domain string,
	entries []crtEntry,
) []string {
	domain =
		strings.ToLower(
			strings.TrimSpace(
				domain,
			),
		)

	domain =
		strings.TrimSuffix(
			domain,
			".",
		)

	seen :=
		make(
			map[string]struct{},
		)

	for _, entry := range entries {

		values :=
			strings.Split(
				entry.NameValue,
				"\n",
			)

		for _, value := range values {

			name :=
				normalizeName(
					value,
				)

			if name == "" {

				continue
			}

			if name == domain {

				continue
			}

			if !strings.HasSuffix(
				name,
				"."+domain,
			) {

				continue
			}

			if strings.Contains(
				name,
				"*",
			) {

				continue
			}

			seen[name] =
				struct{}{}
		}
	}

	result :=
		make(
			[]string,
			0,
			len(seen),
		)

	for name := range seen {

		result =
			append(
				result,
				name,
			)
	}

	sort.Strings(
		result,
	)

	return result
}

func normalizeName(
	value string,
) string {
	value =
		strings.ToLower(
			strings.TrimSpace(
				value,
			),
		)

	value =
		strings.TrimSuffix(
			value,
			".",
		)

	value =
		strings.TrimPrefix(
			value,
			"*.",
		)

	return value
}
