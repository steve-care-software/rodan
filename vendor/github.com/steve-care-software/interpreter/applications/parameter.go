package applications

import "github.com/steve-care-software/interpreter/domain/instructions/parameters"

type parameter struct {
	allParameterIndex   uint
	inputParameterIndex uint
	parameter           parameters.Parameter
}
