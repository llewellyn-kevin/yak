package parser

import (
	"errors"
	"llewellyn-kevin/yak/lexer"
)

type Parser interface {
	Strategy() ParsingStrategy
	Parse() *Program
}

type ParsingStrategy uint8

const (
	RECURSIVE_DESCENT_STRATEGY = iota
)

var first ParsingStrategy = RECURSIVE_DESCENT_STRATEGY
var last ParsingStrategy = RECURSIVE_DESCENT_STRATEGY

type ParserFactory struct {
	mode ParsingStrategy
	err  error
}

func NewParserFactory() *ParserFactory {
	return &ParserFactory{RECURSIVE_DESCENT_STRATEGY, nil}
}

func (p *ParserFactory) Use(s ParsingStrategy) *ParserFactory {
	if s > last || s < first {
		p.err = errors.New("")
		return p
	}
	p.mode = s
	return p
}

func (p ParserFactory) Get(l *lexer.Lexer) (Parser, error) {
	if p.err != nil {
		return nil, p.err
	}

	switch p.mode {
	case RECURSIVE_DESCENT_STRATEGY:
		return newRdParser(l), nil
	default:
		return newRdParser(l), nil
	}
}
