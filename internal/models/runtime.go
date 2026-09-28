package models

type BrowserRuntime struct {
	Executable string

	ProfileDirectory string
}

var Browser = BrowserRuntime{}
