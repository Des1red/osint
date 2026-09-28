package ip

import (
	"errors"
	"net"
	"strings"
)

func validate(input string) error {
	input = strings.TrimSpace(input)

	if input == "" {
		return errors.New("IP address cannot be empty")
	}

	if net.ParseIP(input) == nil {
		return errors.New("invalid IP address")
	}

	return nil
}
