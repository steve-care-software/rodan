package applications

import (
	ast_applications "github.com/steve-care-software/ast/applications"
	"github.com/steve-care-software/ast/domain/grammars"
	"github.com/steve-care-software/ast/domain/trees"
	interpreter_applications "github.com/steve-care-software/interpreter/applications"
	"github.com/steve-care-software/interpreter/domain/programs"
	"github.com/steve-care-software/interpreter/domain/programs/modules"
	query_applications "github.com/steve-care-software/query/applications"
	"github.com/steve-care-software/query/domain/queries"
)

// FetchModulesFn returns the modules
type FetchModulesFn func() (modules.Modules, error)

// NewBuilder creates a new builder instance
func NewBuilder(
	nameBytesToStringFn interpreter_applications.NameBytesToString,
) Builder {
	grammarApp := ast_applications.NewApplication()
	queryApp := query_applications.NewApplication()
	interpreterApp := interpreter_applications.NewApplication(
		nameBytesToStringFn,
	)

	return createBuilder(
		grammarApp,
		queryApp,
		interpreterApp,
	)
}

// Builder represents an application builder
type Builder interface {
	Create() Builder
	WithGrammar(grammar grammars.Grammar) Builder
	WithQuery(query queries.Query) Builder
	WithFetchModulesFn(fetchModulesFn FetchModulesFn) Builder
	Now() (Application, error)
}

// Application represents the rodan application
type Application interface {
	Lex(values []byte) (trees.Tree, error)
	Parse(tree trees.Tree) (programs.Program, []byte, error)
	Interpret(input []interface{}, program programs.Program) ([]interface{}, error)
}
