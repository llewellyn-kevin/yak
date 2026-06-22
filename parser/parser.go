package parser

import (
	"errors"
	"llewellyn-kevin/yak/ast"
	"llewellyn-kevin/yak/lexer"
	"llewellyn-kevin/yak/parser/rd"
	"slices"
)

type ParsingStrategy string

var legalStrategies = []ParsingStrategy{
	rd.RECURSIVE_DESCENT_STRATEGY,
}

type Parser interface {
	Strategy() ParsingStrategy
	Parse() *ast.Program
}

type ParserFactory struct {
	mode ParsingStrategy
	err  error
}

func NewParserFactory() *ParserFactory {
	return &ParserFactory{rd.RECURSIVE_DESCENT_STRATEGY, nil}
}

func (p *ParserFactory) Use(s ParsingStrategy) *ParserFactory {
	if ok := slices.Contains(legalStrategies, s); !ok {
		p.err = errors.New("Could not find parsing strategy in atlas")
	}
	p.mode = s
	return p
}

func (p ParserFactory) Get(l *lexer.Lexer) (*rd.RdParser, error) {
	if p.err != nil {
		return &rd.RdParser{}, p.err
	}

	switch p.mode {
	case rd.RECURSIVE_DESCENT_STRATEGY:
		return rd.NewRdParser(l), nil
	default:
		return rd.NewRdParser(l), nil
	}
}
