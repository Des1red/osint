package ip

import (
	"osint/internal/engines/ip/geo"
	"osint/internal/engines/ip/history"
	"osint/internal/engines/ip/network"
	"osint/internal/engines/ip/rdap"
	"osint/internal/engines/ip/reputation"
	"osint/internal/logger"
	"osint/internal/models"
)

func IpEngine() (
	IPResult,
	error,
) {
	logger.HeaderStart(
		"IP Enrichment",
	)

	defer logger.HeaderEnd(
		"IP Enrichment",
	)

	ip :=
		models.ScopeInput.IpAddress

	var result IPResult

	logger.Info(
		"Validating IP address...",
	)

	err :=
		validate(
			ip,
		)

	if err != nil {

		return result,
			err
	}

	logger.Info(
		"IP address valid.",
	)

	//
	// RDAP failure is not fatal to the IP
	// engine. Other collectors can still
	// produce useful information.
	//
	logger.Info(
		"Collecting RDAP information...",
	)

	rdapResult,
		err :=
		rdap.Rdap(
			ip,
		)

	if err != nil {

		logger.LogError(
			"RDAP failed",
			err.Error(),
		)

	} else {

		result.RDAP =
			rdapResult

		logger.Info(
			"RDAP collection complete.",
		)
	}

	//
	// Geolocation.
	//
	logger.Info(
		"Collecting geolocation information...",
	)

	geoResult,
		err :=
		geo.Geo(
			ip,
		)

	if err != nil {

		logger.LogError(
			"Geolocation failed",
			err.Error(),
		)

	} else {

		result.Geo =
			geoResult

		logger.Info(
			"Geolocation collection complete.",
		)
	}

	//
	// Network intelligence.
	//
	logger.Info(
		"Collecting network intelligence...",
	)

	networkResult,
		err :=
		network.Network(
			ip,
		)

	if err != nil {

		logger.LogError(
			"Network failed",
			err.Error(),
		)

	} else {

		result.Network =
			networkResult

		logger.Info(
			"Network intelligence complete.",
		)
	}

	//
	// Historical intelligence.
	//
	logger.Info(
		"Collecting historical intelligence...",
	)

	historyResult,
		err :=
		history.History(
			ip,
		)

	if err != nil {

		logger.LogError(
			"History failed",
			err.Error(),
		)

	} else {

		result.History =
			historyResult

		logger.Info(
			"Historical intelligence complete.",
		)
	}

	//
	// Reputation.
	//
	logger.Info(
		"Collecting reputation information...",
	)

	reputationResult,
		err :=
		reputation.Reputation(
			ip,
		)

	if err != nil {

		logger.LogError(
			"Reputation failed",
			err.Error(),
		)

	} else {

		result.Reputation =
			reputationResult

		logger.Info(
			"Reputation collection complete.",
		)
	}

	return result,
		nil
}
