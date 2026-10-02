package xvfb

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"osint/internal/models"
)

type runtimeState struct {
	PID int `json:"pid"`

	Executable string `json:"executable"`

	StartTime uint64 `json:"start_time"`

	Display string `json:"display"`
}

func PrepareState() error {
	state,
		err :=
		readRuntimeState()

	if err != nil {

		return err
	}

	if state == nil {

		return nil
	}

	matches,
		err :=
		runtimeStateMatches(
			state,
		)

	if err != nil {

		return err
	}

	//
	// Only kill the process when its Linux process
	// identity exactly matches the state we stored.
	//
	if matches {

		err =
			syscall.Kill(
				state.PID,
				syscall.SIGKILL,
			)

		if err != nil &&
			!errors.Is(
				err,
				syscall.ESRCH,
			) {

			return fmt.Errorf(
				"kill stale Xvfb process %d: %w",
				state.PID,
				err,
			)
		}

		if err == nil {

			err =
				waitForProcessExit(
					state.PID,
				)

			if err != nil {

				return err
			}
		}
	}

	return removeRuntimeState()
}

func writeRuntimeState(
	value *exec.Cmd,
	display string,
) error {
	if value == nil ||
		value.Process == nil {

		return fmt.Errorf(
			"Xvfb process is not available",
		)
	}

	pid :=
		value.Process.Pid

	executable,
		startTime,
		err :=
		linuxProcessIdentity(
			pid,
		)

	if err != nil {

		return fmt.Errorf(
			"resolve Xvfb process identity: %w",
			err,
		)
	}

	state :=
		runtimeState{
			PID: pid,

			Executable: executable,

			StartTime: startTime,

			Display: display,
		}

	data,
		err :=
		json.MarshalIndent(
			state,
			"",
			"  ",
		)

	if err != nil {

		return fmt.Errorf(
			"encode Xvfb runtime state: %w",
			err,
		)
	}

	path,
		err :=
		models.XvfbStatePath()

	if err != nil {

		return err
	}

	err =
		os.MkdirAll(
			filepath.Dir(
				path,
			),
			0700,
		)

	if err != nil {

		return fmt.Errorf(
			"create Xvfb state directory: %w",
			err,
		)
	}

	temporary,
		err :=
		os.CreateTemp(
			filepath.Dir(
				path,
			),
			".xvfb-state-*.tmp",
		)

	if err != nil {

		return fmt.Errorf(
			"create temporary Xvfb state file: %w",
			err,
		)
	}

	temporaryPath :=
		temporary.Name()

	defer os.Remove(
		temporaryPath,
	)

	err =
		temporary.Chmod(
			0600,
		)

	if err != nil {

		temporary.Close()

		return fmt.Errorf(
			"set Xvfb state permissions: %w",
			err,
		)
	}

	_,
		err =
		temporary.Write(
			data,
		)

	if err != nil {

		temporary.Close()

		return fmt.Errorf(
			"write Xvfb runtime state: %w",
			err,
		)
	}

	err =
		temporary.Sync()

	if err != nil {

		temporary.Close()

		return fmt.Errorf(
			"sync Xvfb runtime state: %w",
			err,
		)
	}

	err =
		temporary.Close()

	if err != nil {

		return fmt.Errorf(
			"close Xvfb runtime state: %w",
			err,
		)
	}

	err =
		os.Rename(
			temporaryPath,
			path,
		)

	if err != nil {

		return fmt.Errorf(
			"commit Xvfb runtime state: %w",
			err,
		)
	}

	return nil
}

func readRuntimeState() (
	*runtimeState,
	error,
) {
	path,
		err :=
		models.XvfbStatePath()

	if err != nil {

		return nil,
			err
	}

	data,
		err :=
		os.ReadFile(
			path,
		)

	if os.IsNotExist(
		err,
	) {

		return nil,
			nil
	}

	if err != nil {

		return nil,
			fmt.Errorf(
				"read Xvfb runtime state: %w",
				err,
			)
	}

	var state runtimeState

	err =
		json.Unmarshal(
			data,
			&state,
		)

	if err != nil {

		return nil,
			fmt.Errorf(
				"decode Xvfb runtime state: %w",
				err,
			)
	}

	if state.PID <= 0 {

		return nil,
			fmt.Errorf(
				"invalid Xvfb process ID in runtime state",
			)
	}

	if strings.TrimSpace(
		state.Executable,
	) == "" {

		return nil,
			fmt.Errorf(
				"missing Xvfb executable in runtime state",
			)
	}

	if state.StartTime == 0 {

		return nil,
			fmt.Errorf(
				"missing Xvfb process start time in runtime state",
			)
	}

	return &state,
		nil
}

func removeRuntimeState() error {
	path,
		err :=
		models.XvfbStatePath()

	if err != nil {

		return err
	}

	err =
		os.Remove(
			path,
		)

	if err == nil ||
		os.IsNotExist(
			err,
		) {

		return nil
	}

	return fmt.Errorf(
		"remove Xvfb runtime state: %w",
		err,
	)
}

func runtimeStateMatches(
	state *runtimeState,
) (
	bool,
	error,
) {
	if state == nil {

		return false,
			nil
	}

	executable,
		startTime,
		err :=
		linuxProcessIdentity(
			state.PID,
		)

	if err != nil {

		if errors.Is(
			err,
			os.ErrNotExist,
		) {

			return false,
				nil
		}

		return false,
			err
	}

	if executable !=
		state.Executable {

		return false,
			nil
	}

	if startTime !=
		state.StartTime {

		return false,
			nil
	}

	return true,
		nil
}

func linuxProcessIdentity(
	pid int,
) (
	string,
	uint64,
	error,
) {
	executablePath :=
		fmt.Sprintf(
			"/proc/%d/exe",
			pid,
		)

	executable,
		err :=
		os.Readlink(
			executablePath,
		)

	if err != nil {

		return "",
			0,
			err
	}

	statPath :=
		fmt.Sprintf(
			"/proc/%d/stat",
			pid,
		)

	data,
		err :=
		os.ReadFile(
			statPath,
		)

	if err != nil {

		return "",
			0,
			err
	}

	stat :=
		string(
			data,
		)

	//
	// Field 2 (comm) is enclosed in parentheses and
	// may itself contain spaces, so normal Fields()
	// parsing cannot begin until after its closing
	// parenthesis.
	//
	end :=
		strings.LastIndex(
			stat,
			")",
		)

	if end == -1 {

		return "",
			0,
			fmt.Errorf(
				"invalid /proc/%d/stat format",
				pid,
			)
	}

	fields :=
		strings.Fields(
			stat[end+1:],
		)

	//
	// After removing PID and comm:
	//
	// fields[0] = field 3 (state)
	// fields[19] = field 22 (starttime)
	//
	if len(fields) <= 19 {

		return "",
			0,
			fmt.Errorf(
				"missing process start time in /proc/%d/stat",
				pid,
			)
	}

	startTime,
		err :=
		strconv.ParseUint(
			fields[19],
			10,
			64,
		)

	if err != nil {

		return "",
			0,
			fmt.Errorf(
				"parse process start time: %w",
				err,
			)
	}

	return executable,
		startTime,
		nil
}

func waitForProcessExit(
	pid int,
) error {
	deadline :=
		time.Now().Add(
			shutdownTimeout,
		)

	processDirectory :=
		fmt.Sprintf(
			"/proc/%d",
			pid,
		)

	for time.Now().Before(
		deadline,
	) {

		_,
			err :=
			os.Stat(
				processDirectory,
			)

		if os.IsNotExist(
			err,
		) {

			return nil
		}

		if err != nil {

			return fmt.Errorf(
				"check stale Xvfb process %d: %w",
				pid,
				err,
			)
		}

		time.Sleep(
			25 * time.Millisecond,
		)
	}

	return fmt.Errorf(
		"stale Xvfb process %d did not exit",
		pid,
	)
}
