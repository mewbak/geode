package ast

import (
	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/geode/pkg/util/log"
)

var typeOperators = []string{"*", "?"}

func validTypeInfoTokens(t lexer.Token) bool {
	for _, op := range typeOperators {
		if op == t.Value {
			return true
		}
	}
	return false

}

func (p *Parser) atType() bool {
	if !p.token.Is(lexer.TokType) {
		return false
	}

	offset := 1
	for validTypeInfoTokens(p.Peek(offset)) {
		offset++
	}

	if p.Peek(offset).Type == lexer.TokIdent {
		return true
	}

	return false
}

// parseType returns a

func (p *Parser) parseType() (t TypeNode) {
	p.requires(lexer.TokType)

	t.Name, _ = p.parseName()

	t.Modifiers = make([]TypeModifier, 0)
	// p.Next()

	for {

		if p.token.Is(lexer.TokQuestionMark) {
			if t.Unknown {
				log.Fatal("Multiple Unknown Type operators for %q used.\n", t.Name)
			}

			t.Unknown = true
			t.Modifiers = append(t.Modifiers, ModifierUnknown)
			p.Next()
			continue
		}

		if p.token.Is(lexer.TokOper) && validTypeInfoTokens(p.token) {
			for _, c := range p.token.Value {
				if c == '*' {
					t.PointerLevel++
					t.Modifiers = append(t.Modifiers, ModifierPointer)
				}
			}
			p.Next()
			continue
		}
		// handle slice type definition `T[]` for some T
		// if p.token.Is(lexer.TokLeftBrace) {
		// 	p.Next()
		// 	p.requires(lexer.TokRightBrace)
		// 	t.Modifiers = append(t.Modifiers, ModifierSlice)
		// 	p.Next()
		// 	continue
		// }

		break

	}

	return t
}
