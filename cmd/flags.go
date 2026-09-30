package cmd

import (
	"osint/internal/models"

	"github.com/spf13/pflag"
)

var showHelp bool

func flags() {
	pflag.StringVarP(
		&models.ScopeInput.FullName,
		"name",
		"n",
		"",
		"search information by full name",
	)

	pflag.IntVar(
		&models.ScopeInput.SurnamePosition,
		"surname-position",
		0,
		"position of surname in full name: 1 = first, 2 = last",
	)

	pflag.StringVarP(
		&models.ScopeInput.IpAddress,
		"ip",
		"i",
		"",
		"search information by IP address",
	)

	pflag.StringVarP(
		&models.ScopeInput.Username,
		"username",
		"u",
		"",
		"search information by username",
	)

	pflag.StringVarP(
		&models.ScopeInput.Domain,
		"domain",
		"d",
		"",
		"enumerate subdomains and check for takeover risks",
	)

	pflag.StringVarP(
		&models.ScopeInput.OutputFile,
		"output",
		"o",
		"",
		"file name to save output",
	)

	pflag.BoolVar(
		&models.OutputFlags.Full,
		"full",
		false,
		"show full result details",
	)

	pflag.BoolVarP(
		&showHelp,
		"help",
		"h",
		false,
		"show this help message",
	)

	pflag.BoolVar(
		&models.BootFlags.Install,
		"install",
		false,
		"install osint binary",
	)

	pflag.BoolVar(
		&models.BootFlags.Uninstall,
		"uninstall",
		false,
		"uninstall osint binary and osint files",
	)

	pflag.BoolVar(
		&models.BootFlags.Debug,
		"debug",
		false,
		"enable debug messages",
	)

	pflag.Parse()
}
