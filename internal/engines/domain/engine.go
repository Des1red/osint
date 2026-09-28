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
	domain :=
		normalizeDomain(
			models.ScopeInput.Domain,
		)

	result :=
		DomainResult{
			Domain: domain,
		}

	err :=
		validate(
			domain,
		)

	if err != nil {

		return result,
			err
	}

	//
	// Stage 1:
	//
	// Root DNS information.
	//

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
	}

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
	}

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
	}

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
	}

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
	}

	//
	// Stage 2:
	//
	// Passive subdomain discovery through
	// Certificate Transparency.
	//

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

		//
		// Stage 3:
		//
		// Resolve currently relevant names and
		// inspect CNAME takeover indicators.
		//
		result.Subdomains =
			resolveSubdomains(
				discovery.Names,
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
