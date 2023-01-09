package modules

import (
	"path/filepath"

	"github.com/steve-care-software/ast/applications"
	"github.com/steve-care-software/ast/domain/grammars"
	"github.com/steve-care-software/ast/domain/grammars/cardinalities"
	"github.com/steve-care-software/ast/domain/grammars/values"
	"github.com/steve-care-software/interpreter/domain/programs/modules"
	rodan_grammars "github.com/steve-care-software/rodan/grammars"
	"github.com/steve-care-software/rodan/queries"
	vm_applications "github.com/steve-care-software/vm/applications"
)

const (

	// ModuleList represents a list module
	ModuleList = 0

	// ModuleListFetchElement represents a list fetch element module
	ModuleListFetchElement = 1

	// ModuleFileOpen represents a file open module
	ModuleFileOpen = 2

	// ModuleFileClose represents a write file close module
	ModuleFileClose = 3

	// ModuleFileLock represents a write file lock module
	ModuleFileLock = 4

	// ModuleFileUnLock represents a write file unLock module
	ModuleFileUnLock = 5

	// ModuleFileInfo represents a write file info module
	ModuleFileInfo = 6

	// ModuleFileRead represents a write file read module
	ModuleFileRead = 7

	// ModuleFileWrite represents a write file write module
	ModuleFileWrite = 8

	// ModuleCastToInt represents the castToInt module
	ModuleCastToInt = 9

	// ModuleCastToUint represents the castToUint module
	ModuleCastToUint = 10

	// ModuleCastToBool represents the castToBool module
	ModuleCastToBool = 11

	// ModuleCastToFloat32 represents the castToFloat32 module
	ModuleCastToFloat32 = 12

	// ModuleCastToFloat64 represents the castToFloat32 module
	ModuleCastToFloat64 = 13

	// ModuleASTValue represents the astValue module
	ModuleASTValue = 14

	// ModuleASTCardinality represents the astCardinality module
	ModuleASTCardinality = 15

	// ModuleASTElement represents the astElement module
	ModuleASTElement = 16

	// ModuleASTContainer represents the astContainer module
	ModuleASTContainer = 17

	// ModuleASTLine represents the astLine module
	ModuleASTLine = 18

	// ModuleASTBlock represents the astBlock module
	ModuleASTBlock = 19

	// ModuleASTSuite represents the astSuite module
	ModuleASTSuite = 20

	// ModuleASTSuites represents the astSuites module
	ModuleASTSuites = 21

	// ModuleASTToken represents the astToken module
	ModuleASTToken = 22

	// ModuleASTEverything represents the astEverything module
	ModuleASTEverything = 23

	// ModuleASTInstance represents the astInstance module
	ModuleASTInstance = 24

	// ModuleASTExternal represents the astExternal module
	ModuleASTExternal = 25

	// ModuleASTChannelCondition represents the astChannelCondition module
	ModuleASTChannelCondition = 26

	// ModuleASTChannel represents the astChannel module
	ModuleASTChannel = 27

	// ModuleASTChannels represents the astChannels module
	ModuleASTChannels = 28

	// ModuleAST represents the ast module
	ModuleAST = 29

	// ModuleASTExecute represents the astExecute module
	ModuleASTExecute = 30

	// ModuleVMLex represents a vm lex module
	ModuleVMLex = 31

	// ModuleVMParse represents a vm parse module
	ModuleVMParse = 32

	// ModuleVMInterpret represents a vm interpret module
	ModuleVMInterpret = 33

	// ModuleVMLexParseThenInterpret represents a vm lex, parse then interpreter module
	ModuleVMLexParseThenInterpret = 34

	// ModuleVMLexParseInterpretThenReturnSingle represents a vm lex, parse, interpreter then return single module
	ModuleVMLexParseInterpretThenReturnSingle = 35
)

// NewApplication creates a new virtual machine application
func NewApplication(modulesFn vm_applications.FetchModulesFn) vm_applications.Application {
	vmAppBuilder := vm_applications.NewBuilder(func(name []byte) string {
		return string(name)
	})

	grammar := rodan_grammars.NewGrammar()
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

// NewVMModulesFuncs creates a new vm modules funcs
func NewVMModulesFuncs(
	basePath string,
	chunkSize uint,
) vm_applications.FetchModulesFn {
	vmApp := NewApplication(newFetchModuleFunc(basePath, chunkSize))
	vmModulesFn := createVM(vmApp).Execute()
	baseModulesFn := newModulesFuncs(basePath, chunkSize)

	// create the funcs list:
	allModulesFuncs := map[uint]modules.ExecuteFn{}
	for idx, fn := range vmModulesFn {
		allModulesFuncs[idx] = fn
	}

	for idx, fn := range baseModulesFn {
		allModulesFuncs[idx] = fn
	}

	modulesIns := newModules(allModulesFuncs)
	return func() (modules.Modules, error) {
		return modulesIns, nil
	}
}

func newFetchModuleFunc(
	basePath string,
	chunkSize uint,
) vm_applications.FetchModulesFn {
	fns := newModulesFuncs(basePath, chunkSize)
	modulesIns := newModules(fns)
	return func() (modules.Modules, error) {
		return modulesIns, nil
	}
}

func newModulesFuncs(
	basePath string,
	chunkSize uint,
) map[uint]modules.ExecuteFn {
	// create the containers module funcs:
	containersFnsMap := createContainers().Execute()

	// create the cast module funcs:
	castFnsMap := createCast().Execute()

	// create the file module funcs:
	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		panic(err)
	}

	fileFnsMap := createFile(absBasePath, chunkSize).Execute()

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
	for idx, fn := range containersFnsMap {
		moduleFuncs[idx] = fn
	}

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
