package modules

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestApplication_Success(t *testing.T) {
	// read the script:
	script, err := ioutil.ReadFile("./../scripts/index.rodan")
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	basePath := "./../scripts"
	virtualMachine := NewApplication(newVMModulesFuncs(basePath, 200))
	tree, err := virtualMachine.Lex(script)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	program, remaining, err := virtualMachine.Parse(tree)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	if len(remaining) > 0 {
		t.Errorf("the script was expected to NOT contain remaining data: %s", remaining)
		return
	}

	output, err := virtualMachine.Interpret([]interface{}{}, program)
	if err != nil {
		t.Errorf("the error was expected to be nil, error returned: %s", err.Error())
		return
	}

	fmt.Printf("\n%s\n", output)
}