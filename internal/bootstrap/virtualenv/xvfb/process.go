package xvfb

import (
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"
)

func newCommand(
	executable string,
	writePipe *os.File,
) *exec.Cmd {
	value :=
		exec.Command(
			executable,

			"-displayfd",
			"3",

			"-screen",
			"0",
			"1920x1080x24",

			"-nolisten",
			"tcp",
		)

	value.ExtraFiles =
		[]*os.File{
			writePipe,
		}

	value.Stdout =
		io.Discard

	value.Stderr =
		io.Discard

	//
	// Linux only.
	//
	// If OSINT is killed unexpectedly, make sure
	// its Xvfb child does not remain running.
	//
	value.SysProcAttr =
		&syscall.SysProcAttr{
			Pdeathsig: syscall.SIGKILL,
		}

	return value
}

func startProcess(
	value *exec.Cmd,
) (
	chan error,
	error,
) {
	err :=
		value.Start()

	if err != nil {

		return nil,
			err
	}

	done :=
		make(
			chan error,
			1,
		)

	go func() {

		done <- value.Wait()

		close(
			done,
		)
	}()

	return done,
		nil
}

func stopProcess(
	value *exec.Cmd,
	done chan error,
) {
	if value == nil {

		return
	}

	if value.Process != nil {

		_ =
			value.Process.Kill()
	}

	if done == nil {

		return
	}

	timer :=
		time.NewTimer(
			shutdownTimeout,
		)

	defer timer.Stop()

	select {

	case <-done:

	case <-timer.C:
	}
}
