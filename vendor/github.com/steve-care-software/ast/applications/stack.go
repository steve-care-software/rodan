package applications

import "github.com/steve-care-software/ast/domain/grammars"

type stack struct {
	token grammars.Token
	lines map[int][]byte
}
