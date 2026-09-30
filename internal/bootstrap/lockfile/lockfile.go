package lockfile

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"osint/internal/models"
)

type Lock struct {
	file *os.File
}

func Acquire() (
	*Lock,
	error,
) {
	appDirectory,
		err :=
		models.AppDir()

	if err != nil {

		return nil,
			err
	}

	err =
		os.MkdirAll(
			appDirectory,
			0700,
		)

	if err != nil {

		return nil,
			fmt.Errorf(
				"failed to create application directory: %w",
				err,
			)
	}

	path,
		err :=
		models.LockFilePath()

	if err != nil {

		return nil,
			err
	}

	file,
		err :=
		os.OpenFile(
			path,
			os.O_CREATE|
				os.O_RDWR,
			0600,
		)

	if err != nil {

		return nil,
			fmt.Errorf(
				"failed to open lock file: %w",
				err,
			)
	}

	err =
		syscall.Flock(
			int(
				file.Fd(),
			),
			syscall.LOCK_EX|
				syscall.LOCK_NB,
		)

	if err != nil {

		pid :=
			lockPID(
				file,
			)

		file.Close()

		if errors.Is(
			err,
			syscall.EWOULDBLOCK,
		) ||
			errors.Is(
				err,
				syscall.EAGAIN,
			) {

			if pid != "" {

				return nil,
					fmt.Errorf(
						"osint is already running with PID %s",
						pid,
					)
			}

			return nil,
				fmt.Errorf(
					"osint is already running",
				)
		}

		return nil,
			fmt.Errorf(
				"failed to acquire process lock: %w",
				err,
			)
	}

	err =
		writePID(
			file,
		)

	if err != nil {

		syscall.Flock(
			int(
				file.Fd(),
			),
			syscall.LOCK_UN,
		)

		file.Close()

		return nil,
			err
	}

	return &Lock{
			file: file,
		},
		nil
}

func (
	lock *Lock,
) Release() error {
	if lock == nil ||
		lock.file == nil {

		return nil
	}

	err :=
		syscall.Flock(
			int(
				lock.file.Fd(),
			),
			syscall.LOCK_UN,
		)

	closeErr :=
		lock.file.Close()

	lock.file =
		nil

	if err != nil {

		return fmt.Errorf(
			"failed to release process lock: %w",
			err,
		)
	}

	if closeErr != nil {

		return fmt.Errorf(
			"failed to close lock file: %w",
			closeErr,
		)
	}

	return nil
}

func writePID(
	file *os.File,
) error {
	err :=
		file.Truncate(
			0,
		)

	if err != nil {

		return fmt.Errorf(
			"failed to clear lock file: %w",
			err,
		)
	}

	_,
		err =
		file.Seek(
			0,
			0,
		)

	if err != nil {

		return fmt.Errorf(
			"failed to seek lock file: %w",
			err,
		)
	}

	pid :=
		strconv.Itoa(
			os.Getpid(),
		)

	_,
		err =
		file.WriteString(
			pid +
				"\n",
		)

	if err != nil {

		return fmt.Errorf(
			"failed to write lock PID: %w",
			err,
		)
	}

	return nil
}

func lockPID(
	file *os.File,
) string {
	if file == nil {

		return ""
	}

	_,
		err :=
		file.Seek(
			0,
			0,
		)

	if err != nil {

		return ""
	}

	data,
		err :=
		os.ReadFile(
			file.Name(),
		)

	if err != nil {

		return ""
	}

	return strings.TrimSpace(
		string(
			data,
		),
	)
}
