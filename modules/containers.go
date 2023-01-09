package modules

import (
	"errors"
	"fmt"

	"github.com/steve-care-software/interpreter/domain/programs/modules"
)

type containers struct {
}

func createContainers() *containers {
	out := containers{}
	return &out
}

// Execute executes the application
func (app *containers) Execute() map[uint]modules.ExecuteFn {
	list := app.list()
	fetchElement := app.fetchElement()
	return map[uint]modules.ExecuteFn{
		ModuleList:             list,
		ModuleListFetchElement: fetchElement,
	}
}

func (app *containers) fetchElement() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		if index, ok := input[0].(uint); ok {
			if value, ok := input[1].([]interface{}); ok {
				amount := uint(len(value))
				if index+1 < amount {
					str := fmt.Sprintf("the element at index %d could not be fetched because the list only contains %d elements", index, amount)
					return nil, errors.New(str)
				}

				return value[index], nil
			}

			str := fmt.Sprintf("the input at index %d was expected to contain a list", 1)
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the input at index %d was expected to contain a uint", 0)
		return nil, errors.New(str)
	}
}

func (app *containers) list() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		findValueAtIndex := func(index uint, list map[uint]interface{}) (interface{}, error) {
			for listIndex, element := range list {
				if listIndex != index {
					continue
				}

				return element, nil
			}

			str := fmt.Sprintf("the value at index: %d could not be found in the provided list", index)
			return nil, errors.New(str)
		}

		values := []interface{}{}
		for {
			index := uint(len(values))
			element, err := findValueAtIndex(index, input)
			if err != nil {
				break
			}

			values = append(values, element)
			delete(input, index)
		}

		return values, nil
	}
}
