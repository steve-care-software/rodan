package queries

import (
	"github.com/steve-care-software/interpreter/domain/instructions"
	"github.com/steve-care-software/interpreter/domain/instructions/applications"
	"github.com/steve-care-software/interpreter/domain/instructions/attachments"
	"github.com/steve-care-software/interpreter/domain/instructions/modules"
	"github.com/steve-care-software/interpreter/domain/instructions/parameters"
	"github.com/steve-care-software/query/domain/queries"
)

// NewQuery creates a new query instance
func NewQuery() queries.Query {
	builder := queries.NewBuilder()
	queryFnBuilder := queries.NewQueryFnBuilder()
	tokenBuilder := queries.NewTokenBuilder()
	elementBuilder := queries.NewElementBuilder()
	insideBuilder := queries.NewInsideBuilder()
	fetchersBuilder := queries.NewFetchersBuilder()
	fetcherBuilder := queries.NewFetcherBuilder()
	contentFnBuilder := queries.NewContentFnBuilder()
	instructionsBuilder := instructions.NewBuilder()
	instructionBuilder := instructions.NewInstructionBuilder()
	instructionApplicationBuilder := applications.NewBuilder()
	instructionParameterBuilder := parameters.NewBuilder()
	instructionAttachmentBuilder := attachments.NewBuilder()
	instructionAttachmentVariableBuilder := attachments.NewVariableBuilder()
	instructionAssignmentBuilder := instructions.NewAssignmentBuilder()
	instructionValueBuilder := instructions.NewValueBuilder()
	instructionModuleBuilder := modules.NewBuilder()
	queryIns := createQuery(
		builder,
		queryFnBuilder,
		tokenBuilder,
		elementBuilder,
		insideBuilder,
		fetchersBuilder,
		fetcherBuilder,
		contentFnBuilder,
		instructionsBuilder,
		instructionBuilder,
		instructionApplicationBuilder,
		instructionParameterBuilder,
		instructionAttachmentBuilder,
		instructionAttachmentVariableBuilder,
		instructionAssignmentBuilder,
		instructionValueBuilder,
		instructionModuleBuilder,
	)

	ins, err := queryIns.Execute()
	if err != nil {
		panic(err)
	}

	return ins
}
