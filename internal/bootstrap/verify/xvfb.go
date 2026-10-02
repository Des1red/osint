package verify

import (
	"fmt"
	"os/exec"
)

func VerifyXvfb() error {
	_,
		err :=
		XvfbExecutable()

	return err
}

func XvfbExecutable() (
	string,
	error,
) {
	path,
		err :=
		exec.LookPath(
			"Xvfb",
		)

	if err != nil {

		return "",
			fmt.Errorf(
				"Xvfb is not installed",
			)
	}

	return path,
		nil
}
