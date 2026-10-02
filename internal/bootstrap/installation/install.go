package installation

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"osint/internal/bootstrap/installation/dependencies"
	"osint/internal/bootstrap/installation/privilege"
	"osint/internal/models"
)

const installPath = models.InstallPath

const binName = models.BinName

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
		installBinary()

	if err != nil {

		return fmt.Errorf(
			"install binary: %w",
			err,
		)
	}

	return nil
}

func installBinary() error {
	cmd :=
		exec.Command(
			"go",
			"build",
			"-o",
			binName,
			".",
		)

	cmd.Stdout =
		os.Stdout

	cmd.Stderr =
		os.Stderr

	err :=
		cmd.Run()

	if err != nil {

		return fmt.Errorf(
			"build failed: %w",
			err,
		)
	}

	defer func() {

		err :=
			os.Remove(
				binName,
			)

		if err != nil &&
			!os.IsNotExist(
				err,
			) {

			fmt.Println(
				"warning: could not remove temporary binary:",
				err,
			)
		}
	}()

	command,
		err :=
		privilege.Command(
			"install",
			"-m",
			"0755",
			binName,
			installPath,
		)

	if err != nil {

		return err
	}

	command.Stdout =
		os.Stdout

	command.Stderr =
		os.Stderr

	command.Stdin =
		os.Stdin

	err =
		command.Run()

	if err != nil {

		return fmt.Errorf(
			"install binary to %s: %w",
			installPath,
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
