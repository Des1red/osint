package domain

import (
	subdomains "osint/internal/engines/domain/subdomain"
	"osint/internal/logger"
	"osint/internal/models"
)

func DomainEngine() (
	DomainResult,
	error,
) {
	logger.HeaderStart(
		"Domain",
	)

	defer logger.HeaderEnd(
		"Domain",
	)

	domain :=
		normalizeDomain(
			models.ScopeInput.Domain,
		)

	result :=
		DomainResult{
			Domain: domain,
		}

	logger.Info(
		"Validating domain...",
	)

	err :=
		validate(
			domain,
		)

	if err != nil {

		return result,
			err
	}

	logger.Info(
		"Domain valid.",
	)

	//
	// Stage 1:
	//
	// Root DNS information.
	//
	logger.Info(
		"Looking up DNS addresses...",
	)

	ipv4,
		ipv6,
		err :=
		lookupAddresses(
			domain,
		)

	if err != nil {

		logger.LogError(
			"DNS address lookup failed",
			err.Error(),
		)

	} else {

		result.A =
			ipv4

		result.AAAA =
			ipv6

		logger.Info(
			"DNS address lookup complete.",
		)
	}

	logger.Info(
		"Looking up CNAME records...",
	)

	cname,
		err :=
		lookupCNAME(
			domain,
		)

	if err != nil {

		logger.LogError(
			"DNS CNAME lookup failed",
			err.Error(),
		)

	} else {

		result.CNAME =
			cname

		logger.Info(
			"CNAME lookup complete.",
		)
	}

	logger.Info(
		"Looking up MX records...",
	)

	mx,
		err :=
		lookupMX(
			domain,
		)

	if err != nil {

		logger.LogError(
			"DNS MX lookup failed",
			err.Error(),
		)

	} else {

		result.MX =
			mx

		logger.Info(
			"MX lookup complete.",
		)
	}

	logger.Info(
		"Looking up NS records...",
	)

	ns,
		err :=
		lookupNS(
			domain,
		)

	if err != nil {

		logger.LogError(
			"DNS NS lookup failed",
			err.Error(),
		)

	} else {

		result.NS =
			ns

		logger.Info(
			"NS lookup complete.",
		)
	}

	logger.Info(
		"Looking up TXT records...",
	)

	txt,
		err :=
		lookupTXT(
			domain,
		)

	if err != nil {

		logger.LogError(
			"DNS TXT lookup failed",
			err.Error(),
		)

	} else {

		result.TXT =
			txt

		logger.Info(
			"TXT lookup complete.",
		)
	}

	//
	// Stage 2:
	//
	// Passive subdomain discovery through
	// Certificate Transparency.
	//
	logger.Info(
		"Discovering passive subdomains...",
	)

	discovery,
		err :=
		subdomains.Discover(
			domain,
		)

	if err != nil {

		logger.LogError(
			"Subdomain discovery failed",
			err.Error(),
		)

	} else {

		result.SubdomainSource =
			discovery.Source

		result.SubdomainCandidates =
			discovery.Total

		result.SubdomainChecked =
			len(
				discovery.Names,
			)

		logger.Info(
			"Subdomain discovery complete.",
		)

		//
		// Stage 3:
		//
		// Resolve currently relevant names and
		// inspect CNAME takeover indicators.
		//
		logger.Info(
			"Resolving discovered subdomains...",
		)

		result.Subdomains =
			resolveSubdomains(
				discovery.Names,
			)

		logger.Info(
			"Subdomain resolution complete.",
		)
	}

	result.Found =
		len(result.A) > 0 ||
			len(result.AAAA) > 0 ||
			result.CNAME != "" ||
			len(result.MX) > 0 ||
			len(result.NS) > 0 ||
			len(result.TXT) > 0 ||
			len(result.Subdomains) > 0

	return result,
		nil
}
