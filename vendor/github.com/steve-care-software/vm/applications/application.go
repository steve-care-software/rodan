package applications

import (
	"errors"

	ast_applications "github.com/steve-care-software/ast/applications"
	"github.com/steve-care-software/ast/domain/grammars"
	"github.com/steve-care-software/ast/domain/trees"
	interpreter_applications "github.com/steve-care-software/interpreter/applications"
	"github.com/steve-care-software/interpreter/domain/instructions"
	"github.com/steve-care-software/interpreter/domain/programs"
	query_applications "github.com/steve-care-software/query/applications"
	"github.com/steve-care-software/query/domain/queries"
)

type application struct {
	astApplication         ast_applications.Application
	queryApplication       query_applications.Application
	interpreterApplication interpreter_applications.Application
	grammar                grammars.Grammar
	query                  queries.Query
	fetchModulesFn         FetchModulesFn
}

func createApplication(
	astApplication ast_applications.Application,
	queryApplication query_applications.Application,
	interpreterApplication interpreter_applications.Application,
	grammar grammars.Grammar,
	query queries.Query,
	fetchModulesFn FetchModulesFn,
) Application {
	out := application{
		astApplication:         astApplication,
		queryApplication:       queryApplication,
		interpreterApplication: interpreterApplication,
		grammar:                grammar,
		query:                  query,
		fetchModulesFn:         fetchModulesFn,
	}

	return &out
}

// Lex lexes values into an AST
func (app *application) Lex(values []byte) (trees.Tree, error) {
	return app.astApplication.Execute(app.grammar, values)
}

// Parse parses an AST into a program
func (app *application) Parse(tree trees.Tree) (programs.Program, []byte, error) {
	ins, isValid, remaining, err := app.queryApplication.Execute(app.query, tree)
	if err != nil {
		return nil, nil, err
	}

	if !isValid {
		return nil, remaining, errors.New("the provided AST is not compatible with the VM's Query instance and therefore cannot be parsed by it")
	}

	if castedInstructions, ok := ins.(instructions.Instructions); ok {
		modules, err := app.fetchModulesFn()
		if err != nil {
			return nil, remaining, err
		}

		program, err := app.interpreterApplication.Compile(modules, castedInstructions)
		if err != nil {
			return nil, remaining, err
		}

		return program, remaining, nil
	}

	return nil, remaining, errors.New("the VM's Query instance was expected to return instructions")
}

// Interpret interprets a program with input and returns its output
func (app *application) Interpret(input []interface{}, program programs.Program) ([]interface{}, error) {
	return app.interpreterApplication.Execute(input, program)
}
