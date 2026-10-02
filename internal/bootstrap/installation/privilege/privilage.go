package privilege

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func EnsureUserInvocation() error {
	if os.Geteuid() != 0 {

		return nil
	}

	sudoUser :=
		strings.TrimSpace(
			os.Getenv(
				"SUDO_USER",
			),
		)

	if sudoUser == "" {

		//
		// The actual root user is allowed.
		//
		return nil
	}

	return fmt.Errorf(
		"do not run OSINT installation with sudo; run the command normally and OSINT will elevate only the required system operations",
	)
}

func Command(
	executable string,
	args ...string,
) (
	*exec.Cmd,
	error,
) {
	if os.Geteuid() == 0 {

		return exec.Command(
				executable,
				args...,
			),
			nil
	}

	_,
		err :=
		exec.LookPath(
			"sudo",
		)

	if err != nil {

		return nil,
			fmt.Errorf(
				"sudo is required for this operation: %w",
				err,
			)
	}

	elevatedArgs :=
		[]string{
			executable,
		}

	elevatedArgs =
		append(
			elevatedArgs,
			args...,
		)

	return exec.Command(
			"sudo",
			elevatedArgs...,
		),
		nil
}
