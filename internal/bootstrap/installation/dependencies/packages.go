package dependencies

import (
	"fmt"
	"os"
	"os/exec"

	"osint/internal/bootstrap/installation/privilege"
)

const (
	packageManagerDNF = "dnf"

	packageManagerAPT = "apt-get"
)

func packageManager() (
	string,
	error,
) {
	candidates :=
		[]string{
			packageManagerDNF,
			packageManagerAPT,
		}

	for _, candidate := range candidates {

		_,
			err :=
			exec.LookPath(
				candidate,
			)

		if err == nil {

			return candidate,
				nil
		}
	}

	return "",
		fmt.Errorf(
			"no supported package manager found",
		)
}

func installPackage(
	manager string,
	packageName string,
) error {
	var args []string

	switch manager {

	case packageManagerDNF:

		args =
			[]string{
				"-y",
				"install",
				packageName,
			}

	case packageManagerAPT:

		args =
			[]string{
				"install",
				"-y",
				packageName,
			}

	default:

		return fmt.Errorf(
			"unsupported package manager %q",
			manager,
		)
	}

	command,
		err :=
		privilege.Command(
			manager,
			args...,
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
			"package manager failed: %w",
			err,
		)
	}

	return nil
}

func removePackage(
	manager string,
	packageName string,
) error {
	var args []string

	switch manager {

	case packageManagerDNF:

		args =
			[]string{
				"-y",
				"remove",
				packageName,
			}

	case packageManagerAPT:

		args =
			[]string{
				"remove",
				"-y",
				packageName,
			}

	default:

		return fmt.Errorf(
			"unsupported package manager %q",
			manager,
		)
	}

	command,
		err :=
		privilege.Command(
			manager,
			args...,
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
			"package manager failed: %w",
			err,
		)
	}

	return nil
}
