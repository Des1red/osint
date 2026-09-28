package models

type Scope struct {
	FullName string

	IpAddress string

	Username string

	Domain string

	OutputFile string
}

var ScopeInput = Scope{}

type Engines struct {
	FullName bool

	Ip bool

	Username bool

	Domain bool
}

// Engines requested by the user.
var ActiveEngines = Engines{}

// Engines which actually completed.
//
// Information/output should use this state,
// not ActiveEngines.
//
// An engine can therefore fail without
// preventing other requested engines from
// continuing.
var CompletedEngines = Engines{}
