package virtualenv

import (
	"osint/internal/bootstrap/virtualenv/xvfb"
)

func PrepareState() error {
	return xvfb.PrepareState()
}

func Xvfb() error {
	return xvfb.Start()
}

func Close() error {
	return xvfb.Close()
}
