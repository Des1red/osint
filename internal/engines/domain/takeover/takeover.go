package takeover

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"
)

const lookupTimeout = 5 * time.Second

type Result struct {
	Potential bool

	Provider string

	Target string

	Reason string
}

type provider struct {
	Name string

	Suffixes []string
}

var providers = []provider{
	{
		Name: "GitHub Pages",

		Suffixes: []string{
			"github.io",
		},
	},
	{
		Name: "Heroku",

		Suffixes: []string{
			"herokuapp.com",
		},
	},
	{
		Name: "Netlify",

		Suffixes: []string{
			"netlify.app",
			"netlify.com",
		},
	},
	{
		Name: "Fastly",

		Suffixes: []string{
			"fastly.net",
			"global.ssl.fastly.net",
		},
	},
	{
		Name: "Shopify",

		Suffixes: []string{
			"shops.myshopify.com",
		},
	},
	{
		Name: "Zendesk",

		Suffixes: []string{
			"zendesk.com",
		},
	},
	{
		Name: "Cargo Collective",

		Suffixes: []string{
			"cargocollective.com",
		},
	},
	{
		Name: "Tumblr",

		Suffixes: []string{
			"tumblr.com",
		},
	},
	{
		Name: "AWS CloudFront",

		Suffixes: []string{
			"cloudfront.net",
		},
	},
	{
		Name: "Cloudflare Pages",

		Suffixes: []string{
			"pages.dev",
		},
	},
}

func Check(
	cname string,
) Result {
	target :=
		normalizeHost(
			cname,
		)

	if target == "" {

		return Result{}
	}

	providerName :=
		providerFor(
			target,
		)

	if providerName == "" {

		return Result{
			Target: target,
		}
	}

	result :=
		Result{
			Provider: providerName,

			Target: target,
		}

	notFound,
		err :=
		targetNotFound(
			target,
		)

	if err != nil {

		return result
	}

	if !notFound {

		return result
	}

	result.Potential =
		true

	result.Reason =
		"known third-party CNAME target returns NXDOMAIN; manual validation required"

	return result
}

func providerFor(
	target string,
) string {
	target =
		normalizeHost(
			target,
		)

	for _, candidate := range providers {

		for _, suffix := range candidate.Suffixes {

			suffix =
				normalizeHost(
					suffix,
				)

			if target == suffix ||
				strings.HasSuffix(
					target,
					"."+suffix,
				) {

				return candidate.Name
			}
		}
	}

	return ""
}

func targetNotFound(
	target string,
) (
	bool,
	error,
) {
	ctx,
		cancel :=
		context.WithTimeout(
			context.Background(),
			lookupTimeout,
		)

	defer cancel()

	_,
		err :=
		net.DefaultResolver.LookupHost(
			ctx,
			target,
		)

	if err == nil {

		return false,
			nil
	}

	var dnsError *net.DNSError

	if !errors.As(
		err,
		&dnsError,
	) {

		return false,
			err
	}

	if dnsError.IsNotFound {

		return true,
			nil
	}

	//
	// Timeout, SERVFAIL and other temporary
	// resolver failures are not enough to
	// claim a potential takeover.
	//
	return false,
		err
}

func normalizeHost(
	value string,
) string {
	return strings.ToLower(
		strings.TrimSuffix(
			strings.TrimSpace(
				value,
			),
			".",
		),
	)
}
