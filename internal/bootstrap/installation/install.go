package installation

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	instll "github.com/Des1red/goinstall/cmd"

	"osint/internal/bootstrap/installation/dependencies"
	"osint/internal/bootstrap/installation/privilege"
	"osint/internal/models"
)

func Install() error {
	err :=
		privilege.EnsureUserInvocation()

	if err != nil {
		return err
	}

	installDependencies,
		err :=
		confirmDependencyInstallation()

	if err != nil {
		return fmt.Errorf(
			"dependency confirmation failed: %w",
			err,
		)
	}

	if installDependencies {
		err =
			dependencies.Install()

		if err != nil {
			return fmt.Errorf(
				"install dependencies: %w",
				err,
			)
		}
	} else {
		fmt.Println(
			"dependency installation skipped",
		)

		fmt.Println(
			"required system dependencies must be installed manually",
		)
	}

	err =
		instll.SetBinaryName(
			models.BinName,
		)

	if err != nil {
		return fmt.Errorf(
			"set binary name: %w",
			err,
		)
	}

	err =
		instll.Install(
			true,
			false,
		)

	if err != nil {
		return fmt.Errorf(
			"install binary: %w",
			err,
		)
	}

	return nil
}

func confirmDependencyInstallation() (
	bool,
	error,
) {
	reader :=
		bufio.NewReader(
			os.Stdin,
		)

	for {
		fmt.Print(
			"Install missing system dependencies? [Y/n]: ",
		)

		value,
			err :=
			reader.ReadString(
				'\n',
			)

		if err != nil {
			return false,
				err
		}

		value =
			strings.ToLower(
				strings.TrimSpace(
					value,
				),
			)

		switch value {

		case "",
			"y",
			"yes":

			return true,
				nil

		case "n",
			"no":

			return false,
				nil

		default:

			fmt.Println(
				"please answer yes or no",
			)
		}
	}
}
