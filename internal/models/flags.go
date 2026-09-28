package models

type BootStrap struct {
	Install bool

	Uninstall bool

	Debug bool
}

type OutputOptions struct {
	Full bool
}

var BootFlags = BootStrap{}

var OutputFlags = OutputOptions{}
