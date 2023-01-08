package modules

import (
	"errors"
	"fmt"

	"github.com/steve-care-software/interpreter/domain/programs/modules"
)

type file struct {
	basePath string
}

func createFile(
	basePath string,
) *file {
	out := file{
		basePath: basePath,
	}

	return &out
}

// Execute executes the application
func (app *file) Execute() map[uint]modules.ExecuteFn {
	fileRead := app.fileRead()
	fileWrite := app.fileWrite()
	return map[uint]modules.ExecuteFn{
		ModuleFileRead:  fileRead,
		ModuleFileWrite: fileWrite,
	}
}

func (app *file) fileRead() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		/*if name, ok := input[0].(uint); ok {
			if value, ok := input[1].(map[uint]interface{}); ok {
				return value[name], nil
			}

			str := fmt.Sprintf("the input at index %d was expected to contain a map with uint keynames", 1)
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the input at index %d was expected to contain a uint", 0)
		return nil, errors.New(str)*/
		return nil, nil
	}
}

func (app *file) fileWrite() modules.ExecuteFn {
	return func(input map[uint]interface{}) (interface{}, error) {
		/*if name, ok := input[0].(string); ok {
			if value, ok := input[1].(map[string]interface{}); ok {
				return value[name], nil
			}

			str := fmt.Sprintf("the input at index %d was expected to contain a map with string keynames", 1)
			return nil, errors.New(str)
		}

		str := fmt.Sprintf("the input at index %d was expected to contain a string", 0)
		return nil, errors.New(str)*/
		return nil, nil
	}
}

func (app *file) containerListFetchValue() modules.ExecuteFn {
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

func (app *file) containerList() modules.ExecuteFn {
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
