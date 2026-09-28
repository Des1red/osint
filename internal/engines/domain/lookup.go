package domain

import (
	"context"
	"errors"
	"net"
	"sort"
	"strings"
	"time"
)

const dnsLookupTimeout = 5 * time.Second

func lookupAddresses(
	domain string,
) (
	[]string,
	[]string,
	error,
) {
	ctx,
		cancel :=
		context.WithTimeout(
			context.Background(),
			dnsLookupTimeout,
		)

	defer cancel()

	values,
		err :=
		net.DefaultResolver.LookupIP(
			ctx,
			"ip",
			domain,
		)

	if err != nil {

		if dnsNotFound(
			err,
		) {

			return nil,
				nil,
				nil
		}

		return nil,
			nil,
			err
	}

	var ipv4 []string

	var ipv6 []string

	seen4 :=
		make(
			map[string]struct{},
		)

	seen6 :=
		make(
			map[string]struct{},
		)

	for _, value := range values {

		if value.To4() != nil {

			address :=
				value.String()

			if _, exists :=
				seen4[address]; exists {

				continue
			}

			seen4[address] =
				struct{}{}

			ipv4 =
				append(
					ipv4,
					address,
				)

			continue
		}

		address :=
			value.String()

		if _, exists :=
			seen6[address]; exists {

			continue
		}

		seen6[address] =
			struct{}{}

		ipv6 =
			append(
				ipv6,
				address,
			)
	}

	sort.Strings(
		ipv4,
	)

	sort.Strings(
		ipv6,
	)

	return ipv4,
		ipv6,
		nil
}

func lookupCNAME(
	domain string,
) (
	string,
	error,
) {
	ctx,
		cancel :=
		context.WithTimeout(
			context.Background(),
			dnsLookupTimeout,
		)

	defer cancel()

	value,
		err :=
		net.DefaultResolver.LookupCNAME(
			ctx,
			domain,
		)

	if err != nil {

		if dnsNotFound(
			err,
		) {

			return "",
				nil
		}

		return "",
			err
	}

	value =
		normalizeDNSName(
			value,
		)

	if strings.EqualFold(
		value,
		domain,
	) {

		return "",
			nil
	}

	return value,
		nil
}

func lookupMX(
	domain string,
) (
	[]MXRecord,
	error,
) {
	ctx,
		cancel :=
		context.WithTimeout(
			context.Background(),
			dnsLookupTimeout,
		)

	defer cancel()

	values,
		err :=
		net.DefaultResolver.LookupMX(
			ctx,
			domain,
		)

	if err != nil {

		if dnsNotFound(
			err,
		) {

			return nil,
				nil
		}

		return nil,
			err
	}

	result :=
		make(
			[]MXRecord,
			0,
			len(values),
		)

	for _, value := range values {

		if value == nil {

			continue
		}

		result =
			append(
				result,
				MXRecord{
					Preference: value.Pref,

					Host: normalizeDNSName(
						value.Host,
					),
				},
			)
	}

	sort.Slice(
		result,
		func(
			left int,
			right int,
		) bool {
			if result[left].Preference ==
				result[right].Preference {

				return result[left].Host <
					result[right].Host
			}

			return result[left].Preference <
				result[right].Preference
		},
	)

	return result,
		nil
}

func lookupNS(
	domain string,
) (
	[]string,
	error,
) {
	ctx,
		cancel :=
		context.WithTimeout(
			context.Background(),
			dnsLookupTimeout,
		)

	defer cancel()

	values,
		err :=
		net.DefaultResolver.LookupNS(
			ctx,
			domain,
		)

	if err != nil {

		if dnsNotFound(
			err,
		) {

			return nil,
				nil
		}

		return nil,
			err
	}

	result :=
		make(
			[]string,
			0,
			len(values),
		)

	seen :=
		make(
			map[string]struct{},
		)

	for _, value := range values {

		if value == nil {

			continue
		}

		host :=
			normalizeDNSName(
				value.Host,
			)

		if host == "" {

			continue
		}

		if _, exists :=
			seen[host]; exists {

			continue
		}

		seen[host] =
			struct{}{}

		result =
			append(
				result,
				host,
			)
	}

	sort.Strings(
		result,
	)

	return result,
		nil
}

func lookupTXT(
	domain string,
) (
	[]string,
	error,
) {
	ctx,
		cancel :=
		context.WithTimeout(
			context.Background(),
			dnsLookupTimeout,
		)

	defer cancel()

	values,
		err :=
		net.DefaultResolver.LookupTXT(
			ctx,
			domain,
		)

	if err != nil {

		if dnsNotFound(
			err,
		) {

			return nil,
				nil
		}

		return nil,
			err
	}

	result :=
		make(
			[]string,
			0,
			len(values),
		)

	seen :=
		make(
			map[string]struct{},
		)

	for _, value := range values {

		value =
			strings.TrimSpace(
				value,
			)

		if value == "" {

			continue
		}

		if _, exists :=
			seen[value]; exists {

			continue
		}

		seen[value] =
			struct{}{}

		result =
			append(
				result,
				value,
			)
	}

	sort.Strings(
		result,
	)

	return result,
		nil
}

func normalizeDNSName(
	value string,
) string {
	return strings.TrimSuffix(
		strings.TrimSpace(
			value,
		),
		".",
	)
}

func dnsNotFound(
	err error,
) bool {
	var dnsError *net.DNSError

	if !errors.As(
		err,
		&dnsError,
	) {

		return false
	}

	return dnsError.IsNotFound
}
