package cmd

import (
	"fmt"
	"time"
)

func Run() {
	start :=
		time.Now()

	defer func() {

		fmt.Printf(
			"\nRuntime: %s\n",
			time.Since(
				start,
			).Round(
				time.Millisecond,
			),
		)
	}()

	release :=
		preboot()

	defer release()

	checkflags()

	boot()
	defer cleanup()

	initiate()
}
