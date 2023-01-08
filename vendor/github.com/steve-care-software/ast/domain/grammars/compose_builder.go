package grammars

import (
	"errors"
)

type composeBuilder struct {
	list []ComposeElement
}

func createComposeBuilder() ComposeBuilder {
	out := composeBuilder{
		list: nil,
	}

	return &out
}

// Create initializes the builder
func (app *composeBuilder) Create() ComposeBuilder {
	return createComposeBuilder()
}

// WithList adds a list to the builder
func (app *composeBuilder) WithList(list []ComposeElement) ComposeBuilder {
	app.list = list
	return app
}

// Now builds a new Compose instance
func (app *composeBuilder) Now() (Compose, error) {
	if app.list != nil && len(app.list) <= 0 {
		app.list = nil
	}

	if app.list == nil {
		return nil, errors.New("there must be at least 1 ComposeElement in order to build a Compose instance")
	}

	return createCompose(app.list), nil
}
