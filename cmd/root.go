package cmd

func Run() {
	release :=
		preboot()

	defer release()

	checkflags()

	boot()

	initiate()
}
