package engines

import (
	"strings"

	"osint/internal/engines/domain"
	"osint/internal/engines/fullname"
	"osint/internal/engines/ip"
	"osint/internal/engines/username"
	"osint/internal/knowledge"
	"osint/internal/logger"
	"osint/internal/models"
)

func Manager() {
	activate()

	controller()
}

func activate() {
	module :=
		models.ScopeInput

	if strings.TrimSpace(
		module.FullName,
	) != "" {

		models.ActiveEngines.FullName =
			true
	}

	if strings.TrimSpace(
		module.IpAddress,
	) != "" {

		models.ActiveEngines.Ip =
			true
	}

	if strings.TrimSpace(
		module.Username,
	) != "" {

		models.ActiveEngines.Username =
			true
	}

	if strings.TrimSpace(
		module.Domain,
	) != "" {

		models.ActiveEngines.Domain =
			true
	}
}

func controller() {
	engines :=
		models.ActiveEngines

	if engines.FullName {

		result,
			err :=
			fullname.FullNameEngine()

		if err != nil {

			logger.LogError(
				"Full name engine failed",
				err.Error(),
			)

		} else {

			knowledge.Data.FullName =
				result

			models.CompletedEngines.FullName =
				true
		}
	}

	if engines.Ip {

		result,
			err :=
			ip.IpEngine()

		if err != nil {

			logger.LogError(
				"IP engine failed",
				err.Error(),
			)

		} else {

			knowledge.Data.IP =
				result

			models.CompletedEngines.Ip =
				true
		}
	}

	if engines.Username {

		result,
			err :=
			username.UsernameEngine()

		if err != nil {

			logger.LogError(
				"Username engine failed",
				err.Error(),
			)

		} else {

			knowledge.Data.Username =
				result

			models.CompletedEngines.Username =
				true
		}
	}

	if engines.Domain {

		result,
			err :=
			domain.DomainEngine()

		if err != nil {

			logger.LogError(
				"Domain engine failed",
				err.Error(),
			)

		} else {

			knowledge.Data.Domain =
				result

			models.CompletedEngines.Domain =
				true
		}
	}
}
