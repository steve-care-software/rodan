package queries

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/steve-care-software/ast/domain/trees"
	"github.com/steve-care-software/interpreter/domain/instructions"
	"github.com/steve-care-software/interpreter/domain/instructions/applications"
	"github.com/steve-care-software/interpreter/domain/instructions/attachments"
	"github.com/steve-care-software/interpreter/domain/instructions/modules"
	"github.com/steve-care-software/interpreter/domain/instructions/parameters"
	"github.com/steve-care-software/query/domain/queries"
)

type query struct {
	builder                              queries.Builder
	queryFnBuilder                       queries.QueryFnBuilder
	tokenBuilder                         queries.TokenBuilder
	elementBuilder                       queries.ElementBuilder
	insideBuilder                        queries.InsideBuilder
	fetchersBuilder                      queries.FetchersBuilder
	fetcherBuilder                       queries.FetcherBuilder
	contentFnBuilder                     queries.ContentFnBuilder
	instructionsBuilder                  instructions.Builder
	instructionBuilder                   instructions.InstructionBuilder
	instructionApplicationBuilder        applications.Builder
	instructionParameterBuilder          parameters.Builder
	instructionAttachmentBuilder         attachments.Builder
	instructionAttachmentVariableBuilder attachments.VariableBuilder
	instructionAssignmentBuilder         instructions.AssignmentBuilder
	instructionValueBuilder              instructions.ValueBuilder
	instructionModuleBuilder             modules.Builder
}

func createQuery(
	builder queries.Builder,
	queryFnBuilder queries.QueryFnBuilder,
	tokenBuilder queries.TokenBuilder,
	elementBuilder queries.ElementBuilder,
	insideBuilder queries.InsideBuilder,
	fetchersBuilder queries.FetchersBuilder,
	fetcherBuilder queries.FetcherBuilder,
	contentFnBuilder queries.ContentFnBuilder,
	instructionsBuilder instructions.Builder,
	instructionBuilder instructions.InstructionBuilder,
	instructionApplicationBuilder applications.Builder,
	instructionParameterBuilder parameters.Builder,
	instructionAttachmentBuilder attachments.Builder,
	instructionAttachmentVariableBuilder attachments.VariableBuilder,
	instructionAssignmentBuilder instructions.AssignmentBuilder,
	instructionValueBuilder instructions.ValueBuilder,
	instructionModuleBuilder modules.Builder,
) *query {
	out := query{
		builder:                              builder,
		queryFnBuilder:                       queryFnBuilder,
		tokenBuilder:                         tokenBuilder,
		elementBuilder:                       elementBuilder,
		instructionsBuilder:                  instructionsBuilder,
		insideBuilder:                        insideBuilder,
		fetchersBuilder:                      fetchersBuilder,
		fetcherBuilder:                       fetcherBuilder,
		contentFnBuilder:                     contentFnBuilder,
		instructionBuilder:                   instructionBuilder,
		instructionApplicationBuilder:        instructionApplicationBuilder,
		instructionParameterBuilder:          instructionParameterBuilder,
		instructionAttachmentBuilder:         instructionAttachmentBuilder,
		instructionAttachmentVariableBuilder: instructionAttachmentVariableBuilder,
		instructionAssignmentBuilder:         instructionAssignmentBuilder,
		instructionValueBuilder:              instructionValueBuilder,
		instructionModuleBuilder:             instructionModuleBuilder,
	}

	return &out
}

// Execute returns the query
func (app *query) Execute() (queries.Query, error) {
	return app.instructions(), nil
}

func (app *query) instructions() queries.Query {
	return app.queryWithMultiFn(
		app.token(
			"instructions",
			app.element("instruction", 0),
		),
		app.insideWithQueries([]queries.Query{
			app.moduleDeclaration(),
			app.applicationDeclaration(),
			app.parameter(),
			app.execute("instruction"),
			app.assignment(),
			app.attachment(),
		}),
		func(instances []interface{}) (interface{}, bool, error) {
			list := []instructions.Instruction{}
			for idx, oneIns := range instances {
				if casted, ok := oneIns.(instructions.Instruction); ok {
					list = append(list, casted)
					continue
				}

				str := fmt.Sprintf("the instruction (index: %d) could not be properly casted", idx)
				return nil, false, errors.New(str)
			}

			ins, err := app.instructionsBuilder.Create().
				WithList(list).
				Now()

			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) attachment() queries.Query {
	return app.queryWithMultiFn(
		app.tokenWithContentIndex(
			"instruction",
			app.element("attachment", 0),
			0,
		),
		app.insideWithQueries([]queries.Query{
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"attachment",
					app.element("variableReference", 0),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"attachment",
					app.element("attachmentTarget", 0),
					0,
				),
				app.insideWithQuery(
					app.queryWithSingleFn(
						app.tokenWithContentIndex(
							"attachmentTarget",
							app.element("number", 0),
							0,
						),
						app.fetchAllContentInside(),
						func(instance interface{}) (interface{}, bool, error) {
							return instance.([]byte), true, nil
						},
					),
				),
				func(instance interface{}) (interface{}, bool, error) {
					return instance, true, nil
				},
			),
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"attachment",
					app.element("variableReference", 1),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
		}),
		func(instances []interface{}) (interface{}, bool, error) {
			if len(instances) != 3 {
				str := fmt.Sprintf("%d elements were expected, %d returned", 3, len(instances))
				return nil, false, errors.New(str)
			}

			target, err := strconv.Atoi(string(instances[1].([]byte)))
			if err != nil {
				return nil, false, err
			}

			variable, err := app.instructionAttachmentVariableBuilder.Create().
				WithCurrent(instances[0].([]byte)).
				WithTarget(uint(target)).
				Now()

			if err != nil {
				return nil, false, err
			}

			attachment, err := app.instructionAttachmentBuilder.Create().
				WithVariable(variable).
				WithApplication(instances[2].([]byte)).
				Now()

			if err != nil {
				return nil, false, err
			}

			ins, err := app.instructionBuilder.Create().
				WithAttachment(attachment).
				Now()

			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) assignment() queries.Query {
	return app.queryWithSingleFn(
		app.token(
			"instruction",
			app.element("assignment", 0),
		),
		app.insideWithQueries([]queries.Query{
			app.variableAssignment(),
			app.executionAssignment(),
			app.instructionsAssignment(),
			app.constantAssignment(),
		}),
		func(instance interface{}) (interface{}, bool, error) {
			if casted, ok := instance.(instructions.Assignment); ok {
				ins, err := app.instructionBuilder.Create().
					WithAssignment(casted).
					Now()

				if err != nil {
					return nil, false, err
				}

				return ins, true, nil
			}

			return nil, false, errors.New("the assignment could not be casted properly")
		},
	)
}

func (app *query) constantAssignment() queries.Query {
	return app.queryWithMultiFn(
		app.tokenWithContentIndex(
			"assignment",
			app.element("constantAssignment", 0),
			0,
		),
		app.insideWithQueries([]queries.Query{
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"constantAssignment",
					app.element("variableReference", 0),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"constantAssignment",
					app.element("everythingExceptEndOfLine", 0),
					0,
				),
				app.fetchAllContentInside(),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
		}),
		func(instances []interface{}) (interface{}, bool, error) {
			if len(instances) != 2 {
				str := fmt.Sprintf("%d elements were expected, %d returned", 2, len(instances))
				return nil, false, errors.New(str)
			}

			value, err := app.instructionValueBuilder.Create().
				WithConstant(instances[1].([]byte)).
				Now()

			if err != nil {
				return nil, false, err
			}

			ins, err := app.instructionAssignmentBuilder.Create().
				WithVariable(instances[0].([]byte)).
				WithValue(value).
				Now()

			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) instructionsAssignment() queries.Query {
	return app.queryWithMultiFn(
		app.tokenWithContentIndex(
			"assignment",
			app.element("instructionsAssignment", 0),
			0,
		),
		app.insideWithQueries([]queries.Query{
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"instructionsAssignment",
					app.element("variableReference", 0),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"instructionsAssignment",
					app.element("instructions", 0),
					0,
				),
				app.insideWithRecursive("instructions"),
				func(instance interface{}) (interface{}, bool, error) {
					if casted, ok := instance.(instructions.Instructions); ok {
						return casted, true, nil
					}

					return nil, false, errors.New("the instance was expected to contain an Instructions instance")
				},
			),
		}),
		func(instances []interface{}) (interface{}, bool, error) {
			if len(instances) != 2 {
				str := fmt.Sprintf("%d elements were expected, %d returned", 2, len(instances))
				return nil, false, errors.New(str)
			}

			builder := app.instructionAssignmentBuilder.Create().
				WithVariable(instances[0].([]byte))

			if casted, ok := instances[1].(instructions.Instructions); ok {
				value, err := app.instructionValueBuilder.Create().
					WithInstructions(casted).
					Now()

				if err != nil {
					return nil, false, err
				}

				builder.WithValue(value)
			}

			ins, err := builder.Now()
			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) executionAssignment() queries.Query {
	return app.queryWithMultiFn(
		app.tokenWithContentIndex(
			"assignment",
			app.element("executionAssignment", 0),
			0,
		),
		app.insideWithQueries([]queries.Query{
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"executionAssignment",
					app.element("variableReference", 0),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
			app.execute("executionAssignment"),
		}),
		func(instances []interface{}) (interface{}, bool, error) {
			if len(instances) != 2 {
				str := fmt.Sprintf("%d elements were expected, %d returned", 2, len(instances))
				return nil, false, errors.New(str)
			}

			builder := app.instructionAssignmentBuilder.Create().
				WithVariable(instances[0].([]byte))

			if casted, ok := instances[1].(instructions.Instruction); ok {
				if !casted.IsExecution() {
					return nil, false, errors.New("the Instruction was expected to contain an Execute name")
				}

				execution := casted.Execution()
				value, err := app.instructionValueBuilder.Create().WithExecution(execution).Now()
				if err != nil {
					return nil, false, err
				}

				builder.WithValue(value)
			}

			ins, err := builder.Now()
			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) variableAssignment() queries.Query {
	return app.queryWithMultiFn(
		app.tokenWithContentIndex(
			"assignment",
			app.element("variableAssignment", 0),
			0,
		),
		app.insideWithQueries([]queries.Query{
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"variableAssignment",
					app.element("variableReference", 0),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"variableAssignment",
					app.element("variableReference", 1),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
		}),
		func(instances []interface{}) (interface{}, bool, error) {
			if len(instances) != 2 {
				str := fmt.Sprintf("%d elements were expected, %d returned", 2, len(instances))
				return nil, false, errors.New(str)
			}

			value, err := app.instructionValueBuilder.Create().WithVariable(instances[1].([]byte)).Now()
			if err != nil {
				return nil, false, err
			}

			ins, err := app.instructionAssignmentBuilder.Create().
				WithVariable(instances[0].([]byte)).
				WithValue(value).
				Now()

			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) execute(tokenName string) queries.Query {
	return app.queryWithSingleFn(
		app.tokenWithContentIndex(
			tokenName,
			app.element("execute", 0),
			0,
		),
		app.insideWithQuery(
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"execute",
					app.element("variableReference", 0),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					ins, err := app.instructionBuilder.Create().
						WithExecution(instance.([]byte)).
						Now()

					if err != nil {
						return nil, false, err
					}

					return ins, true, nil
				},
			),
		),
		func(instance interface{}) (interface{}, bool, error) {
			if casted, ok := instance.(instructions.Instruction); ok {
				return casted, true, nil
			}

			return nil, false, errors.New("the instruction could not be casted properly")
		},
	)
}

func (app *query) parameter() queries.Query {
	return app.queryWithSingleFn(
		app.token(
			"instruction",
			app.element("parameter", 0),
		),
		app.insideWithQueries([]queries.Query{
			app.inputParameter(),
			app.outputParameter(),
		}),
		func(instance interface{}) (interface{}, bool, error) {
			if casted, ok := instance.(parameters.Parameter); ok {
				ins, err := app.instructionBuilder.Create().
					WithParameter(casted).
					Now()

				if err != nil {
					return nil, false, err
				}

				return ins, true, nil
			}

			return nil, false, errors.New("the parameter could not be casted properly")
		},
	)
}

func (app *query) outputParameter() queries.Query {
	return app.queryWithMultiFn(
		app.tokenWithContentIndex(
			"parameter",
			app.element("outputParameter", 0),
			0,
		),
		app.insideWithQuery(
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"outputParameter",
					app.element("variableReference", 0),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
		),
		func(instances []interface{}) (interface{}, bool, error) {
			ins, err := app.instructionParameterBuilder.Create().
				WithName(instances[0].([]byte)).
				Now()

			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) inputParameter() queries.Query {
	return app.queryWithMultiFn(
		app.tokenWithContentIndex(
			"parameter",
			app.element("inputParameter", 0),
			0,
		),
		app.insideWithQuery(
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"inputParameter",
					app.element("variableReference", 0),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
		),
		func(instances []interface{}) (interface{}, bool, error) {
			ins, err := app.instructionParameterBuilder.Create().
				WithName(instances[0].([]byte)).
				IsInput().
				Now()

			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) applicationDeclaration() queries.Query {
	return app.queryWithMultiFn(
		app.tokenWithContentIndex(
			"instruction",
			app.element("applicationDeclaration", 0),
			0,
		),
		app.insideWithQueries([]queries.Query{
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"applicationDeclaration",
					app.element("variableReference", 0),
					0,
				),
				app.insideWithQuery(app.variableReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"applicationDeclaration",
					app.element("moduleReference", 0),
					0,
				),
				app.insideWithQuery(app.moduleReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance.([]byte), true, nil
				},
			),
		}),
		func(instances []interface{}) (interface{}, bool, error) {
			if len(instances) != 2 {
				str := fmt.Sprintf("%d elements were expected, %d returned", 2, len(instances))
				return nil, false, errors.New(str)
			}

			appIns, err := app.instructionApplicationBuilder.Create().
				WithName(instances[0].([]byte)).
				WithModule(instances[1].([]byte)).
				Now()

			if err != nil {
				return nil, false, err
			}

			ins, err := app.instructionBuilder.Create().
				WithApplication(appIns).
				Now()

			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) variableReference() queries.Query {
	return app.queryWithSingleFn(
		app.tokenWithContentIndex(
			"variableReference",
			app.element("name", 0),
			0,
		),
		app.fetchAllContentInside(),
		func(instance interface{}) (interface{}, bool, error) {
			return instance.([]byte), true, nil
		},
	)
}

func (app *query) moduleDeclaration() queries.Query {
	return app.queryWithMultiFn(
		app.tokenWithContentIndex(
			"instruction",
			app.element("moduleDeclaration", 0),
			0,
		),
		app.insideWithQueries([]queries.Query{
			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"moduleDeclaration",
					app.element("moduleReference", 0),
					0,
				),
				app.insideWithQuery(app.moduleReference()),
				func(instance interface{}) (interface{}, bool, error) {
					return instance, true, nil
				},
			),

			app.queryWithSingleFn(
				app.tokenWithContentIndex(
					"moduleDeclaration",
					app.element("moduleIndex", 0),
					0,
				),
				app.insideWithQuery(
					app.queryWithSingleFn(
						app.tokenWithContentIndex(
							"moduleIndex",
							app.element("number", 0),
							0,
						),
						app.fetchAllContentInside(),
						func(instance interface{}) (interface{}, bool, error) {
							return instance.([]byte), true, nil
						},
					),
				),
				func(instance interface{}) (interface{}, bool, error) {
					return instance, true, nil
				},
			),
		}),
		func(instances []interface{}) (interface{}, bool, error) {
			if len(instances) != 2 {
				str := fmt.Sprintf("%d elements were expected, %d returned", 2, len(instances))
				return nil, false, errors.New(str)
			}

			index, err := strconv.Atoi(string(instances[1].([]byte)))
			if err != nil {
				return nil, false, nil
			}

			moduleIns, err := app.instructionModuleBuilder.Create().
				WithName(instances[0].([]byte)).
				WithIndex(uint(index)).
				Now()

			if err != nil {
				return nil, false, err
			}

			ins, err := app.instructionBuilder.Create().
				WithModule(moduleIns).
				Now()

			if err != nil {
				return nil, false, err
			}

			return ins, true, nil
		},
	)
}

func (app *query) moduleReference() queries.Query {
	return app.queryWithSingleFn(
		app.tokenWithContentIndex(
			"moduleReference",
			app.element("name", 0),
			0,
		),
		app.fetchAllContentInside(),
		func(instance interface{}) (interface{}, bool, error) {
			return instance.([]byte), true, nil
		},
	)
}

func (app *query) fetchAllContentInside() queries.Inside {
	return app.insideWithFn(
		app.contentFnWithSingle(
			func(content trees.Content) ([]interface{}, error) {
				return []interface{}{
					content.Bytes(false),
				}, nil
			},
		),
	)
}

func (app *query) queryWithSingleFn(
	token queries.Token,
	inside queries.Inside,
	fn queries.SingleQueryFn,
) queries.Query {
	queryFn, err := app.queryFnBuilder.Create().
		WithSingle(fn).
		Now()

	if err != nil {
		panic(err)
	}

	return app.query(token, inside, queryFn)
}

func (app *query) queryWithMultiFn(
	token queries.Token,
	inside queries.Inside,
	fn queries.MultiQueryFn,
) queries.Query {

	queryFn, err := app.queryFnBuilder.Create().
		WithMulti(fn).
		Now()

	if err != nil {
		panic(err)
	}

	return app.query(token, inside, queryFn)
}

func (app *query) query(
	token queries.Token,
	inside queries.Inside,
	fn queries.QueryFn,
) queries.Query {
	query, err := app.builder.Create().
		WithToken(token).
		WithInside(inside).
		WithFn(fn).
		Now()

	if err != nil {
		panic(err)
	}

	return query
}

func (app *query) insideWithRecursive(recursive string) queries.Inside {
	return app.insideWithRecursives([]string{
		recursive,
	})
}

func (app *query) insideWithRecursives(recursives []string) queries.Inside {
	fetchersList := []queries.Fetcher{}
	for _, oneRecursive := range recursives {
		fetcher, err := app.fetcherBuilder.Create().WithRecursive(oneRecursive).Now()
		if err != nil {
			panic(err)
		}

		fetchersList = append(fetchersList, fetcher)
	}

	fetchers, err := app.fetchersBuilder.Create().WithList(fetchersList).Now()
	if err != nil {
		panic(err)
	}

	inside, err := app.insideBuilder.Create().
		WithFetchers(fetchers).
		Now()

	if err != nil {
		panic(err)
	}

	return inside
}

func (app *query) insideWithQuery(query queries.Query) queries.Inside {
	return app.insideWithQueries([]queries.Query{
		query,
	})
}

func (app *query) insideWithQueries(queriesList []queries.Query) queries.Inside {
	fetchersList := []queries.Fetcher{}
	for _, oneQuery := range queriesList {
		fetcher, err := app.fetcherBuilder.Create().WithQuery(oneQuery).Now()
		if err != nil {
			panic(err)
		}

		fetchersList = append(fetchersList, fetcher)
	}

	fetchers, err := app.fetchersBuilder.Create().WithList(fetchersList).Now()
	if err != nil {
		panic(err)
	}

	inside, err := app.insideBuilder.Create().
		WithFetchers(fetchers).
		Now()

	if err != nil {
		panic(err)
	}

	return inside
}

func (app *query) insideWithFn(fn queries.ContentFn) queries.Inside {
	inside, err := app.insideBuilder.Create().
		WithFn(fn).
		Now()

	if err != nil {
		panic(err)
	}

	return inside
}

func (app *query) contentFnWithSingle(fn queries.SingleContentFn) queries.ContentFn {
	contentFn, err := app.contentFnBuilder.Create().
		WithSingle(fn).
		Now()

	if err != nil {
		panic(err)
	}

	return contentFn
}

func (app *query) token(
	name string,
	element queries.Element,
) queries.Token {
	ins, err := app.tokenBuilder.Create().
		WithName(name).
		WithElement(element).
		WithReverseName("reverse").
		Now()

	if err != nil {
		panic(err)
	}

	return ins
}

func (app *query) tokenWithContentIndex(
	name string,
	element queries.Element,
	contentIndex uint,
) queries.Token {
	ins, err := app.tokenBuilder.Create().
		WithName(name).
		WithElement(element).
		WithReverseName("reverse").
		WithContent(contentIndex).
		Now()

	if err != nil {
		panic(err)
	}

	return ins
}

func (app *query) element(
	name string,
	index uint,
) queries.Element {
	ins, err := app.elementBuilder.Create().
		WithName(name).
		WithIndex(index).
		Now()

	if err != nil {
		panic(err)
	}

	return ins
}
