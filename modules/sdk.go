package modules

import (
	"github.com/steve-care-software/ast/applications"
	"github.com/steve-care-software/ast/domain/grammars"
	"github.com/steve-care-software/ast/domain/grammars/cardinalities"
	"github.com/steve-care-software/ast/domain/grammars/values"
	"github.com/steve-care-software/interpreter/domain/programs/modules"
	vm_applications "github.com/steve-care-software/vm/applications"
)

const (

	// ModuleCastToInt represents the castToInt module
	ModuleCastToInt = 0

	// ModuleCastToUint represents the castToUint module
	ModuleCastToUint = 1

	// ModuleCastToBool represents the castToBool module
	ModuleCastToBool = 2

	// ModuleCastToFloat32 represents the castToFloat32 module
	ModuleCastToFloat32 = 3

	// ModuleCastToFloat64 represents the castToFloat32 module
	ModuleCastToFloat64 = 4

	// ModuleFileRead represents a read file module
	ModuleFileRead = 5

	// ModuleFileWrite represents a write file module
	ModuleFileWrite = 6

	// ModuleASTValue represents the astValue module
	ModuleASTValue = 7

	// ModuleASTCardinality represents the astCardinality module
	ModuleASTCardinality = 8

	// ModuleASTElement represents the astElement module
	ModuleASTElement = 9

	// ModuleASTContainer represents the astContainer module
	ModuleASTContainer = 10

	// ModuleASTLine represents the astLine module
	ModuleASTLine = 11

	// ModuleASTBlock represents the astBlock module
	ModuleASTBlock = 12

	// ModuleASTSuite represents the astSuite module
	ModuleASTSuite = 13

	// ModuleASTSuites represents the astSuites module
	ModuleASTSuites = 14

	// ModuleASTToken represents the astToken module
	ModuleASTToken = 15

	// ModuleASTEverything represents the astEverything module
	ModuleASTEverything = 16

	// ModuleASTInstance represents the astInstance module
	ModuleASTInstance = 17

	// ModuleASTExternal represents the astExternal module
	ModuleASTExternal = 18

	// ModuleASTChannelCondition represents the astChannelCondition module
	ModuleASTChannelCondition = 19

	// ModuleASTChannel represents the astChannel module
	ModuleASTChannel = 20

	// ModuleASTChannels represents the astChannels module
	ModuleASTChannels = 21

	// ModuleAST represents the ast module
	ModuleAST = 22

	// ModuleASTExecute represents the astExecute module
	ModuleASTExecute = 23
)

// NewFetchModuleFunc creates a new fetchModule func
func NewFetchModuleFunc(
	basePath string,
) vm_applications.FetchModulesFn {
	fns := newModulesFuncs(basePath)
	modulesIns := newModules(fns)
	return func() (modules.Modules, error) {
		return modulesIns, nil
	}
}

func newModulesFuncs(
	basePath string,
) map[uint]modules.ExecuteFn {
	// create the cast module funcs:
	castFnsMap := createCast().Execute()

	// create the file module funcs:
	fileFnsMap := createFile(basePath).Execute()

	// create the ast module funcs:
	astApplication := applications.NewApplication()
	grammarBuilder := grammars.NewBuilder()
	grammarChannelsBuilder := grammars.NewChannelsBuilder()
	grammarChannelBuilder := grammars.NewChannelBuilder()
	grammarChannelConditionBuilder := grammars.NewChannelConditionBuilder()
	grammarExternalBuilder := grammars.NewExternalBuilder()
	grammarInstanceBuilder := grammars.NewInstanceBuilder()
	grammarEverythingBuilder := grammars.NewEverythingBuilder()
	grammarTokenBuilder := grammars.NewTokenBuilder()
	grammarSuitesBuilder := grammars.NewSuitesBuilder()
	grammarSuiteBuilder := grammars.NewSuiteBuilder()
	grammarBlockBuilder := grammars.NewBlockBuilder()
	grammarLineBuilder := grammars.NewLineBuilder()
	grammarContainerBuilder := grammars.NewContainerBuilder()
	grammarElementBuilder := grammars.NewElementBuilder()
	grammarCardinalityBuilder := cardinalities.NewBuilder()
	grammarValueBuilder := values.NewBuilder()
	grammarFnsMap := createAST(
		astApplication,
		grammarBuilder,
		grammarChannelsBuilder,
		grammarChannelBuilder,
		grammarChannelConditionBuilder,
		grammarExternalBuilder,
		grammarInstanceBuilder,
		grammarEverythingBuilder,
		grammarTokenBuilder,
		grammarSuitesBuilder,
		grammarSuiteBuilder,
		grammarBlockBuilder,
		grammarLineBuilder,
		grammarContainerBuilder,
		grammarElementBuilder,
		grammarCardinalityBuilder,
		grammarValueBuilder,
	).Execute()

	// create the module funcs list:
	moduleFuncs := map[uint]modules.ExecuteFn{}
	for idx, fn := range castFnsMap {
		moduleFuncs[idx] = fn
	}

	for idx, fn := range fileFnsMap {
		moduleFuncs[idx] = fn
	}

	for idx, fn := range grammarFnsMap {
		moduleFuncs[idx] = fn
	}

	return moduleFuncs
}

func newModules(moduleFuncs map[uint]modules.ExecuteFn) modules.Modules {
	// build the modules list:
	modulesList := []modules.Module{}
	moduleBuilder := modules.NewModuleBuilder()
	for idx, oneFunc := range moduleFuncs {
		ins, err := moduleBuilder.Create().WithIndex(uint(idx)).WithFunc(oneFunc).Now()
		if err != nil {
			panic(err)
		}

		modulesList = append(modulesList, ins)
	}

	modulesIns, err := modules.NewBuilder().Create().WithList(modulesList).Now()
	if err != nil {
		panic(err)
	}

	return modulesIns
}
