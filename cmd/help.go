package cmd

import (
	"fmt"

	"github.com/Des1red/clihelp"
)

func help() {
	description := `A passive OSINT tool for collecting publicly available intelligence about names, IP addresses, usernames, and domains.`

	fmt.Println("OSINT-Master")
	fmt.Println()
	fmt.Println(description)
	fmt.Println()

	fmt.Println("Usage:")
	fmt.Println("  osint [flags]")
	fmt.Println()

	fmt.Println("Targets:")
	clihelp.Print(
		clihelp.F(
			"-n, --name",
			`"Full Name"`,
			"search information by full name",
		),
		clihelp.F(
			"--surname-position",
			`"1|2"`,
			"position of surname in full name: 1 = first, 2 = last (default 2)",
		),
		clihelp.F(
			"-i, --ip",
			`"IP Address"`,
			"search information by IP address",
		),
		clihelp.F(
			"-u, --username",
			`"Username"`,
			"search information by username",
		),
		clihelp.F(
			"-d, --domain",
			`"Domain"`,
			"enumerate subdomains and check for takeover risks",
		),
	)

	fmt.Println()
	fmt.Println("Setup:")
	clihelp.Print(
		clihelp.F(
			"--install",
			"",
			"build and install the osint binary",
		),
		clihelp.F(
			"--uninstall",
			"",
			"remove the osint binary and local cache files",
		),
	)

	fmt.Println()
	fmt.Println("General:")
	clihelp.Print(
		clihelp.F(
			"-h, --help",
			"",
			"show this help message",
		),
	)

	fmt.Println()
	fmt.Println("Output:")
	clihelp.Print(
		clihelp.F(
			"-o, --output",
			`"FileName"`,
			"save results to a file",
		),
		clihelp.F(
			"--full",
			"",
			"show full result details",
		),
	)
}
