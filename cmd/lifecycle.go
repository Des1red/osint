package cmd

import (
	"os"
	"osint/internal/bootstrap"
	"osint/internal/engines"
	"osint/internal/information"
)

func checkflags() {
	flags()
	if showHelp {
		help()
		os.Exit(0)
	}
	if len(os.Args) == 1 {
		help()
		os.Exit(0)
	}
}

func boot() {
	bootstrap.Boot()
}

func initiate() {
	engines.Manager()
	information.Print()

}
