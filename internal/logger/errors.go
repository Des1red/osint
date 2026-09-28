package logger

import "fmt"

func LogError(err string, reason string) {
	fmt.Println("ERROR: " + err + " | Reason: " + reason)
}
