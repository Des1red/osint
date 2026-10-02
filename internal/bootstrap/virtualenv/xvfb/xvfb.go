package xvfb

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Start() error {
	stateMu.Lock()

	defer stateMu.Unlock()

	//
	// Already running.
	//
	if command != nil {

		select {

		case <-processDone:

			restoreDisplay()

			command =
				nil

			processDone =
				nil

		default:

			return nil
		}
	}

	executable,
		err :=
		exec.LookPath(
			"Xvfb",
		)

	if err != nil {

		return fmt.Errorf(
			"resolve Xvfb executable: %w",
			err,
		)
	}

	readPipe,
		writePipe,
		err :=
		os.Pipe()

	if err != nil {

		return fmt.Errorf(
			"create Xvfb display pipe: %w",
			err,
		)
	}

	defer readPipe.Close()

	xvfbCommand :=
		newCommand(
			executable,
			writePipe,
		)

	done,
		err :=
		startProcess(
			xvfbCommand,
		)

	//
	// The child owns its duplicate of the write
	// side now. The parent must close its copy or
	// EOF detection on the read side will not work.
	//
	closeErr :=
		writePipe.Close()

	if err != nil {

		if closeErr != nil {

			return fmt.Errorf(
				"start Xvfb: %w; close display pipe: %v",
				err,
				closeErr,
			)
		}

		return fmt.Errorf(
			"start Xvfb: %w",
			err,
		)
	}

	displayReady :=
		readDisplay(
			readPipe,
		)

	timer :=
		time.NewTimer(
			startupTimeout,
		)

	defer timer.Stop()

	var displayNumber string

	select {

	case result :=
		<-displayReady:

		if result.Err != nil {

			stopProcess(
				xvfbCommand,
				done,
			)

			return fmt.Errorf(
				"read Xvfb display number: %w",
				result.Err,
			)
		}

		displayNumber =
			result.Value

	case processErr :=
		<-done:

		if processErr != nil {

			return fmt.Errorf(
				"Xvfb exited before reporting a display: %w",
				processErr,
			)
		}

		return fmt.Errorf(
			"Xvfb exited before reporting a display",
		)

	case <-timer.C:

		stopProcess(
			xvfbCommand,
			done,
		)

		return fmt.Errorf(
			"Xvfb startup timed out",
		)
	}

	display,
		err :=
		setDisplay(
			displayNumber,
		)

	if err != nil {

		stopProcess(
			xvfbCommand,
			done,
		)

		return err
	}

	//
	// Make sure Xvfb did not terminate immediately
	// after reporting its display.
	//
	select {

	case processErr :=
		<-done:

		restoreDisplay()

		if processErr != nil {

			return fmt.Errorf(
				"Xvfb exited during startup: %w",
				processErr,
			)
		}

		return fmt.Errorf(
			"Xvfb exited during startup",
		)

	default:
	}

	//
	// Persist the Linux process identity only after
	// Xvfb has started successfully and its display
	// has been configured.
	//
	err =
		writeRuntimeState(
			xvfbCommand,
			display,
		)

	if err != nil {

		stopProcess(
			xvfbCommand,
			done,
		)

		restoreDisplay()

		return err
	}

	command =
		xvfbCommand

	processDone =
		done

	return nil
}

func Close() error {
	stateMu.Lock()

	defer stateMu.Unlock()

	if command != nil {

		stopProcess(
			command,
			processDone,
		)
	}

	command =
		nil

	processDone =
		nil

	restoreDisplay()

	err :=
		removeRuntimeState()

	if err != nil {

		return err
	}

	return nil
}
