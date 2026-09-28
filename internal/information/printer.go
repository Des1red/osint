package information

import (
	"fmt"

	"osint/internal/information/domain"
	"osint/internal/information/fullname"
	"osint/internal/information/ip"
	"osint/internal/information/output"
	"osint/internal/information/username"
	"osint/internal/models"
)

func Print() {
	engines :=
		models.CompletedEngines

	if engines.FullName {

		hasData :=
			fullname.FullNamePrinter()

		if hasData {

			err :=
				output.WriteFile(
					"fullname",
					models.ScopeInput.FullName,
					func() {
						fullname.FullNamePrinter()
					},
				)

			if err != nil {

				fmt.Printf(
					"failed to write full name output: %v\n",
					err,
				)
			}
		}
	}

	if engines.Ip {

		hasData :=
			ip.IpPrinter()

		if hasData {

			err :=
				output.WriteFile(
					"ip",
					models.ScopeInput.IpAddress,
					func() {
						ip.IpPrinter()
					},
				)

			if err != nil {

				fmt.Printf(
					"failed to write IP output: %v\n",
					err,
				)
			}
		}
	}

	if engines.Username {

		hasData :=
			username.UsernamePrinter()

		if hasData {

			err :=
				output.WriteFile(
					"username",
					models.ScopeInput.Username,
					func() {
						username.UsernamePrinter()
					},
				)

			if err != nil {

				fmt.Printf(
					"failed to write username output: %v\n",
					err,
				)
			}
		}
	}

	if engines.Domain {

		hasData :=
			domain.DomainPrinter()

		if hasData {

			err :=
				output.WriteFile(
					"domain",
					models.ScopeInput.Domain,
					func() {
						domain.DomainPrinter()
					},
				)

			if err != nil {

				fmt.Printf(
					"failed to write domain output: %v\n",
					err,
				)
			}
		}
	}
}
