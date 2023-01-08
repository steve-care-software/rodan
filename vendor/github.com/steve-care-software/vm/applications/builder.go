package applications

import (
	"errors"

	ast_applications "github.com/steve-care-software/ast/applications"
	"github.com/steve-care-software/ast/domain/grammars"
	interpreter_applications "github.com/steve-care-software/interpreter/applications"
	query_applications "github.com/steve-care-software/query/applications"
	"github.com/steve-care-software/query/domain/queries"
)

type builder struct {
	astApplication         ast_applications.Application
	queryApplication       query_applications.Application
	interpreterApplication interpreter_applications.Application
	grammar                grammars.Grammar
	query                  queries.Query
	fetchModulesFn         FetchModulesFn
}

func createBuilder(
	astApplication ast_applications.Application,
	queryApplication query_applications.Application,
	interpreterApplication interpreter_applications.Application,
) Builder {
	out := builder{
		astApplication:         astApplication,
		queryApplication:       queryApplication,
		interpreterApplication: interpreterApplication,
	}

	return &out
}

// Create initializes the builder
func (app *builder) Create() Builder {
	return createBuilder(
		app.astApplication,
		app.queryApplication,
		app.interpreterApplication,
	)
}

// WithGrammar adds a grammar to the builder
func (app *builder) WithGrammar(grammar grammars.Grammar) Builder {
	app.grammar = grammar
	return app
}

// WithQuery adds a query to the builder
func (app *builder) WithQuery(query queries.Query) Builder {
	app.query = query
	return app
}

// WithFetchModulesFn adds a fetchModulesFn to the builder
func (app *builder) WithFetchModulesFn(fetchModulesFn FetchModulesFn) Builder {
	app.fetchModulesFn = fetchModulesFn
	return app
}

// Now builds a new Application instance
func (app *builder) Now() (Application, error) {
	if app.grammar == nil {
		return nil, errors.New("the grammar is mandatory in order to build an Application instance")
	}

	if app.query == nil {
		return nil, errors.New("the query is mandatory in order to build an Application instance")
	}

	if app.fetchModulesFn == nil {
		return nil, errors.New("the fetchModulesFn is mandatory in order to build an Application instance")
	}

	return createApplication(
		app.astApplication,
		app.queryApplication,
		app.interpreterApplication,
		app.grammar,
		app.query,
		app.fetchModulesFn,
	), nil
}
