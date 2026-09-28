package network

import (
	"net"
	"strings"
)

func reverseDNS(
	ip string,
) ([]string, error) {
	names, err := net.LookupAddr(ip)
	if err != nil {
		return nil, err
	}

	results := make(
		[]string,
		0,
		len(names),
	)

	for _, name := range names {
		name = strings.TrimSuffix(
			name,
			".",
		)

		if name == "" {
			continue
		}

		results = append(
			results,
			name,
		)
	}

	return results, nil
}
