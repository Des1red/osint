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
	ip :=
		models.ScopeInput.IpAddress

	var result IPResult

	err :=
		validate(
			ip,
		)

	if err != nil {

		return result,
			err
	}

	//
	// RDAP failure is not fatal to the IP
	// engine. Other collectors can still
	// produce useful information.
	//
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
	}

	//
	// Geolocation.
	//
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
	}

	//
	// Network intelligence.
	//
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
	}

	//
	// Historical intelligence.
	//
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
	}

	//
	// Reputation.
	//
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
	}

	return result,
		nil
}
