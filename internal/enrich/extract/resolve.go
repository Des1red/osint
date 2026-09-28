package extract

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func resolveURL(
	value string,
) (string, error) {
	resolved, err :=
		resolveWithMethod(
			value,
			http.MethodHead,
		)

	if err == nil {
		return resolved,
			nil
	}

	return resolveWithMethod(
		value,
		http.MethodGet,
	)
}

func resolveWithMethod(
	value string,
	method string,
) (string, error) {
	parsed, err :=
		url.Parse(
			value,
		)

	if err != nil {
		return "",
			err
	}

	if parsed.Scheme != "http" &&
		parsed.Scheme != "https" {

		return "",
			fmt.Errorf(
				"unsupported URL scheme",
			)
	}

	if !publicHost(
		parsed.Hostname(),
	) {
		return "",
			fmt.Errorf(
				"non-public destination",
			)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,

		CheckRedirect: func(
			req *http.Request,
			via []*http.Request,
		) error {
			if len(via) >= 10 {
				return fmt.Errorf(
					"too many redirects",
				)
			}

			if !publicHost(
				req.URL.Hostname(),
			) {
				return fmt.Errorf(
					"redirected to non-public destination",
				)
			}

			return nil
		},
	}

	request, err :=
		http.NewRequest(
			method,
			value,
			nil,
		)

	if err != nil {
		return "",
			err
	}

	request.Header.Set(
		"User-Agent",
		"osint-master",
	)

	response, err :=
		client.Do(
			request,
		)

	if err != nil {
		return "",
			err
	}

	defer response.Body.Close()

	return response.Request.URL.String(),
		nil
}

func publicHost(
	host string,
) bool {
	host =
		strings.TrimSpace(
			host,
		)

	if host == "" {
		return false
	}

	if strings.EqualFold(
		host,
		"localhost",
	) {
		return false
	}

	addresses, err :=
		net.LookupIP(
			host,
		)

	if err != nil {
		return false
	}

	for _, address := range addresses {

		if address.IsLoopback() ||
			address.IsPrivate() ||
			address.IsUnspecified() ||
			address.IsLinkLocalUnicast() ||
			address.IsLinkLocalMulticast() {

			return false
		}
	}

	return true
}
