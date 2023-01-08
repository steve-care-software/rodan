package applications

import (
	"github.com/steve-care-software/rodan/grammars"
	"github.com/steve-care-software/rodan/queries"
	vm_applications "github.com/steve-care-software/vm/applications"
)

// NewApplication creates a new virtual machine application
func NewApplication(modulesFn vm_applications.FetchModulesFn) vm_applications.Application {
	vmAppBuilder := vm_applications.NewBuilder(func(name []byte) string {
		return string(name)
	})

	grammar := grammars.NewGrammar()
	query := queries.NewQuery()
	vmApp, err := vmAppBuilder.Create().
		WithFetchModulesFn(modulesFn).
		WithGrammar(grammar).
		WithQuery(query).
		Now()

	if err != nil {
		panic(err)
	}

	return vmApp
}
