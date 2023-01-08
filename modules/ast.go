package modules

import (
	"errors"
	"fmt"

	"github.com/steve-care-software/ast/applications"
	"github.com/steve-care-software/ast/domain/grammars"
	"github.com/steve-care-software/ast/domain/grammars/cardinalities"
	"github.com/steve-care-software/ast/domain/grammars/values"
	"github.com/steve-care-software/interpreter/domain/programs/modules"
)

type ast struct {
	astApplication          applications.Application
	builder                 grammars.Builder
	channelsBuilder         grammars.ChannelsBuilder
	channelBuilder          grammars.ChannelBuilder
	channelConditionBuilder grammars.ChannelConditionBuilder
	externalBuilder         grammars.ExternalBuilder
	instanceBuilder         grammars.InstanceBuilder
	everythingBuilder       grammars.EverythingBuilder
	tokenBuilder            grammars.TokenBuilder
	suitesBuilder           grammars.SuitesBuilder
	suiteBuilder            grammars.SuiteBuilder
	blockBuilder            grammars.BlockBuilder
	lineBuilder             grammars.LineBuilder
	containerBuilder        grammars.ContainerBuilder
	elementBuilder          grammars.ElementBuilder
	cardinalityBuilder      cardinalities.Builder
	valueBuilder            values.Builder
}

func createAST(
	astApplication applications.Application,
	builder grammars.Builder,
	channelsBuilder grammars.ChannelsBuilder,
	channelBuilder grammars.ChannelBuilder,
	channelConditionBuilder grammars.ChannelConditionBuilder,
	externalBuilder grammars.ExternalBuilder,
	instanceBuilder grammars.InstanceBuilder,
	everythingBuilder grammars.EverythingBuilder,
	tokenBuilder grammars.TokenBuilder,
	suitesBuilder grammars.SuitesBuilder,
	suiteBuilder grammars.SuiteBuilder,
	blockBuilder grammars.BlockBuilder,
	lineBuilder grammars.LineBuilder,
	containerBuilder grammars.ContainerBuilder,
	elementBuilder grammars.ElementBuilder,
	cardinalityBuilder cardinalities.Builder,
	valueBuilder values.Builder,
) *ast {
	out := ast{
		astApplication:          astApplication,
		builder:                 builder,
		channelsBuilder:         channelsBuilder,
		channelBuilder:          channelBuilder,
		channelConditionBuilder: channelConditionBuilder,
		externalBuilder:         externalBuilder,
		instanceBuilder:         instanceBuilder,
		everythingBuilder:       everythingBuilder,
		tokenBuilder:            tokenBuilder,
		suitesBuilder:           suitesBuilder,
		suiteBuilder:            suiteBuilder,
		blockBuilder:            blockBuilder,
		lineBuilder:             lineBuilder,
		containerBuilder:        containerBuilder,
		elementBuilder:          elementBuilder,
		cardinalityBuilder:      cardinalityBuilder,
		valueBuilder:            valueBuilder,
	}

	return &out
}

// Execute executes the application
func (app *ast) Execute() map[uint]modules.ExecuteFn {
	value := app.value()
	cardinality := app.cardinality()
	element := app.element()
	container := app.container()
	line := app.line()
	block := app.block()
	suite := app.suite()
	suites := app.suites()
	token := app.token()
	everything := app.everything()
	instance := app.instance()
	external := app.external()
	channelCondition := app.channelCondition()
	channel := app.channel()
	channels := app.channels()
	grammar := app.grammar()
	execute := app.execute()
	return map[uint]modules.ExecuteFn{
		ModuleASTValue:            value,
		ModuleASTCardinality:      cardinality,
		ModuleASTElement:          element,
		ModuleASTContainer:        container,
		ModuleASTLine:             line,
		ModuleASTBlock:            block,
		ModuleASTSuite:            suite,
		ModuleASTSuites:           suites,
		ModuleASTToken:            token,
		ModuleASTEverything:       everything,
		ModuleASTInstance:         instance,
		ModuleASTExternal:         external,
		ModuleASTChannelCondition: channelCondition,
		ModuleASTChannel:          channel,
		ModuleASTChannels:         channels,
		ModuleAST:                 grammar,
		ModuleASTExecute:          execute,
	}
}

func (app *ast) execute() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if grammar, ok := input[0].(grammars.Grammar); ok {
			if data, ok := input[1].([]byte); ok {
				return app.astApplication.Execute(grammar, data)
			}

			str := fmt.Sprintf("the data was expected to be defined")
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the grammar was expected to be defined")
		return nil, errors.New(str)
	}
}

func (app *ast) grammar() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.builder.Create()
		if root, ok := input[0].(grammars.Token); ok {
			builder.WithRoot(root)
		}

		if channels, ok := input[1].(grammars.Channels); ok {
			builder.WithChannels(channels)
		}

		return builder.Now()
	}
}

func (app *ast) channels() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if channelsList, ok := input[0].([]interface{}); ok {
			list := []grammars.Channel{}
			for index, oneChannel := range channelsList {
				if casted, ok := oneChannel.(grammars.Channel); ok {
					list = append(list, casted)
					continue
				}

				str := fmt.Sprintf("the value at index: %d was expected to be a Channel instance", index)
				return nil, errors.New(str)
			}

			return app.channelsBuilder.Create().WithList(list).Now()
		}

		str := fmt.Sprintf("the channels was expected to be valid and contain a list")
		return nil, errors.New(str)
	}
}

func (app *ast) channel() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.channelBuilder.Create()
		if token, ok := input[0].(grammars.Token); ok {
			builder.WithToken(token)
		}

		if condition, ok := input[1].(grammars.ChannelCondition); ok {
			builder.WithCondition(condition)
		}

		return builder.Now()
	}
}

func (app *ast) channelCondition() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.channelConditionBuilder.Create()
		if previous, ok := input[0].(grammars.Token); ok {
			builder.WithPrevious(previous)
		}

		if next, ok := input[1].(grammars.Token); ok {
			builder.WithNext(next)
		}

		return builder.Now()
	}
}

func (app *ast) external() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.externalBuilder.Create()
		if name, ok := input[0].([]byte); ok {
			builder.WithName(string(name))
		}

		if grammar, ok := input[1].(grammars.Grammar); ok {
			builder.WithGrammar(grammar)
		}

		return builder.Now()
	}
}

func (app *ast) instance() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.instanceBuilder.Create()
		if token, ok := input[0].(grammars.Token); ok {
			builder.WithToken(token)
		}

		if everything, ok := input[1].(grammars.Everything); ok {
			builder.WithEverything(everything)
		}

		return builder.Now()
	}
}

func (app *ast) everything() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.everythingBuilder.Create()
		if name, ok := input[0].([]byte); ok {
			builder.WithName(string(name))
		}

		if exception, ok := input[1].(grammars.Token); ok {
			builder.WithException(exception)
		}

		if escape, ok := input[2].(grammars.Token); ok {
			builder.WithEscape(escape)
		}

		return builder.Now()
	}
}

func (app *ast) token() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.tokenBuilder.Create()
		if name, ok := input[0].([]byte); ok {
			builder.WithName(string(name))
		}

		if block, ok := input[1].(grammars.Block); ok {
			builder.WithBlock(block)
		}

		if suites, ok := input[2].(grammars.Suites); ok {
			builder.WithSuites(suites)
		}

		return builder.Now()
	}
}

func (app *ast) suites() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if suitesList, ok := input[0].([]interface{}); ok {
			list := []grammars.Suite{}
			for index, oneSuite := range suitesList {
				if casted, ok := oneSuite.(grammars.Suite); ok {
					list = append(list, casted)
					continue
				}

				str := fmt.Sprintf("the value at index: %d was expected to be a Suite instance", index)
				return nil, errors.New(str)
			}

			return app.suitesBuilder.Create().WithList(list).Now()
		}

		str := fmt.Sprintf("the suites was expected to be valid and contain a list")
		return nil, errors.New(str)
	}
}

func (app *ast) suite() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.suiteBuilder.Create()
		if valid, ok := input[0].(grammars.Compose); ok {
			builder.WithValid(valid)
		}

		if invalid, ok := input[1].(grammars.Compose); ok {
			builder.WithInvalid(invalid)
		}

		return builder.Now()
	}
}

func (app *ast) block() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if linesList, ok := input[0].([]interface{}); ok {
			list := []grammars.Line{}
			for index, oneLine := range linesList {
				if casted, ok := oneLine.(grammars.Line); ok {
					list = append(list, casted)
					continue
				}

				str := fmt.Sprintf("the value at index: %d was expected to be a Line instance", index)
				return nil, errors.New(str)
			}

			return app.blockBuilder.Create().WithLines(list).Now()
		}

		str := fmt.Sprintf("the lines was expected to be valid and contain a list")
		return nil, errors.New(str)
	}
}

func (app *ast) line() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if elementsList, ok := input[0].([]interface{}); ok {
			list := []grammars.Container{}
			for index, oneContainer := range elementsList {
				if casted, ok := oneContainer.(grammars.Container); ok {
					list = append(list, casted)
					continue
				}

				str := fmt.Sprintf("the value at index: %d was expected to be a Container instance", index)
				return nil, errors.New(str)
			}

			return app.lineBuilder.Create().WithContainers(list).Now()
		}

		str := fmt.Sprintf("the elements was expected to be valid and contain a list")
		return nil, errors.New(str)
	}
}

func (app *ast) container() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.containerBuilder.Create()
		if element, ok := input[0].(grammars.Element); ok {
			builder.WithElement(element)
		}

		if compose, ok := input[1].(grammars.Compose); ok {
			builder.WithCompose(compose)
		}

		return builder.Now()
	}
}

func (app *ast) element() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.elementBuilder.Create()
		if cardinality, ok := input[0].(cardinalities.Cardinality); ok {
			builder.WithCardinality(cardinality)
		}

		if value, ok := input[1].(values.Value); ok {
			builder.WithValue(value)
		}

		if external, ok := input[2].(grammars.External); ok {
			builder.WithExternal(external)
		}

		if instance, ok := input[3].(grammars.Instance); ok {
			builder.WithInstance(instance)
		}

		return builder.Now()
	}
}

func (app *ast) cardinality() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if min, ok := input[0].(uint); ok {
			builder := app.cardinalityBuilder.Create().WithMin(min)
			if max, ok := input[1].(uint); ok {
				if max < 0 {
					return nil, errors.New("the maximum cannot be smaller or equal than 0")
				}

				builder.WithMax(max)
			}

			return builder.Now()
		}

		str := fmt.Sprintf("the name was expected to be valid and contain a string")
		return nil, errors.New(str)
	}
}

func (app *ast) value() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		builder := app.valueBuilder.Create()
		if name, ok := input[1]; ok {
			builder.WithName(fmt.Sprintf("%s", name))
		}

		if number, ok := input[0].(uint); ok {
			if number > 255 {
				return nil, errors.New("the number cannot be bigger than 255")
			}

			builder.WithNumber(byte(number))
		}

		return builder.Now()
	}
}
