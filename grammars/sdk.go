package grammars

import (
	"github.com/steve-care-software/ast/domain/grammars"
	"github.com/steve-care-software/ast/domain/grammars/cardinalities"
	"github.com/steve-care-software/ast/domain/grammars/values"
)

// NewGrammar creates a new grammar instance
func NewGrammar() grammars.Grammar {
	builder := grammars.NewBuilder()
	channelsBuilder := grammars.NewChannelsBuilder()
	channelBuilder := grammars.NewChannelBuilder()
	instanceBuilder := grammars.NewInstanceBuilder()
	everythingBuilder := grammars.NewEverythingBuilder()
	tokensBuilder := grammars.NewTokensBuilder()
	tokenBuilder := grammars.NewTokenBuilder()
	suitesBuilder := grammars.NewSuitesBuilder()
	suiteBuilder := grammars.NewSuiteBuilder()
	blockBuilder := grammars.NewBlockBuilder()
	lineBuilder := grammars.NewLineBuilder()
	containerBuilder := grammars.NewContainerBuilder()
	elementBuilder := grammars.NewElementBuilder()
	valueBuilder := values.NewBuilder()
	cardinalityBuilder := cardinalities.NewBuilder()
	grammarIns := createGrammar(
		builder,
		channelsBuilder,
		channelBuilder,
		instanceBuilder,
		everythingBuilder,
		tokensBuilder,
		tokenBuilder,
		suitesBuilder,
		suiteBuilder,
		blockBuilder,
		lineBuilder,
		containerBuilder,
		elementBuilder,
		valueBuilder,
		cardinalityBuilder,
	)

	ins, err := grammarIns.Execute()
	if err != nil {
		panic(err)
	}

	return ins
}
