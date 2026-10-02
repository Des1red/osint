package xvfb

import (
	"os/exec"
	"sync"
	"time"
)

const startupTimeout = 5 * time.Second

const shutdownTimeout = 2 * time.Second

var stateMu sync.Mutex

var command *exec.Cmd

var processDone chan error

var previousDisplay string

var previousDisplaySet bool

type displayResult struct {
	Value string

	Err error
}
