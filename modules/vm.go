package modules

import (
	"errors"
	"fmt"

	"github.com/steve-care-software/ast/domain/trees"
	"github.com/steve-care-software/interpreter/domain/programs"
	"github.com/steve-care-software/interpreter/domain/programs/modules"
	"github.com/steve-care-software/vm/applications"
)

type vm struct {
	vmApplication applications.Application
}

func createVM(
	vmApplication applications.Application,
) *vm {
	out := vm{
		vmApplication: vmApplication,
	}

	return &out
}

// Execute executes the application
func (app *vm) Execute() map[uint]modules.ExecuteFn {
	lex := app.lex()
	parse := app.parse()
	interpret := app.interpret()
	lexParseThenInterpret := app.lexParseThenInterpret()
	lexParseThenInterpretSingle := app.lexParseThenInterpretSingle()
	return map[uint]modules.ExecuteFn{
		ModuleVMLex:                   lex,
		ModuleVMParse:                 parse,
		ModuleVMInterpret:             interpret,
		ModuleVMLexParseThenInterpret: lexParseThenInterpret,
		ModuleVMLexParseInterpretThenReturnSingle: lexParseThenInterpretSingle,
	}
}

func (app *vm) lex() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if script, ok := input[0].([]byte); ok {
			return app.vmApplication.Lex(script)
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain []byte", 0)
		return nil, errors.New(str)
	}
}

func (app *vm) parse() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if treeIns, ok := input[0].(trees.Tree); ok {
			programIns, remaining, err := app.vmApplication.Parse(treeIns)
			if err != nil {
				return nil, err
			}

			if len(remaining) > 0 {
				return nil, errors.New("the script was expected to NOT contain remaining data")
			}

			return programIns, nil
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain an AST", 0)
		return nil, errors.New(str)
	}
}

func (app *vm) interpret() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if programIns, ok := input[0].(programs.Program); ok {
			params := []interface{}{}
			if inputList, ok := input[1].([]interface{}); ok {
				params = inputList
			}

			return app.vmApplication.Interpret(params, programIns)
		}

		str := fmt.Sprintf("the input at index (%d) was expected to contain a Program", 0)
		return nil, errors.New(str)
	}
}

func (app *vm) lexParseThenInterpret() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		return app.lexParseThenInterpreterInput(input)
	}
}

func (app *vm) lexParseThenInterpretSingle() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		list, err := app.lexParseThenInterpreterInput(input)
		if err != nil {
			return nil, err
		}

		if len(list) <= 0 {
			return nil, errors.New("the lexParseThenInterpretSingle module was expected to return a single element, but the output is empty")
		}

		return list[0], nil
	}
}

func (app *vm) lexParseThenInterpreterInput(input map[uint]interface{}) ([]interface{}, error) {
	if script, ok := input[0].([]byte); ok {
		treeIns, err := app.vmApplication.Lex(script)
		if err != nil {
			return nil, err
		}

		programIns, remaining, err := app.vmApplication.Parse(treeIns)
		if err != nil {
			return nil, err
		}

		if len(remaining) > 0 {
			return nil, errors.New("the script was expected to NOT contain remaining data")
		}

		params := []interface{}{}
		if inputList, ok := input[1].([]interface{}); ok {
			params = inputList
		}

		return app.vmApplication.Interpret(params, programIns)
	}

	str := fmt.Sprintf("the input at index (%d) was expected to contain []byte", 0)
	return nil, errors.New(str)
}
