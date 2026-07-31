package ast

import (
	"fmt"
	"sort"
	"strings"
)

type Program struct {
	FunctionTable map[string]*FunctionStatement
	MainBlock     *Block
}

func (p Program) String() (s string) {
	names := make([]string, 0, len(p.FunctionTable))
	for name := range p.FunctionTable {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		s += p.FunctionTable[name].String() + "\n"
	}
	s += p.MainBlock.String()
	return
}

func (p Program) GetFunction(name string) (*FunctionStatement, error) {
	fn, ok := p.FunctionTable[name]
	if !ok {
		return nil, fmt.Errorf("could not find function with name `%s`", name)
	}
	return fn, nil
}

func (p *Program) AddFunction(fn *FunctionStatement) {
	if p.FunctionTable == nil {
		p.FunctionTable = make(map[string]*FunctionStatement)
	}
	p.FunctionTable[fn.Name] = fn
}

type Block struct {
	Program     *Program
	Id          int
	NestLevel   int
	Scope       string
	IsLoop      bool
	Statements  []Statement
	Expressions []Expression
}

func (Block) isStatement() {}

func (b Block) IsRoot() bool {
	return b.Id == 0
}

func (b Block) IsNamedScope() bool {
	return b.Scope != ""
}

func (b Block) GetBlock(id int) (*Block, error) {
	for _, statement := range b.Statements {
		if block, ok := statement.(*Block); ok {
			if block.Id == id {
				return block, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not find block with id `%d`", id)
}

func (b Block) GetConditional(id int) (*ConditionalStatement, error) {
	for _, statement := range b.Statements {
		if block, ok := statement.(*ConditionalStatement); ok {
			if block.Id == id {
				return block, nil
			}
		}
	}
	return nil, fmt.Errorf("Could not find conditional block with id `%d`", id)
}

func (b Block) String() string {
	indent := strings.Repeat("\t", b.NestLevel*2)
	bodyIndent := strings.Repeat("\t", b.NestLevel*2+1)
	contentIndent := strings.Repeat("\t", b.NestLevel*2+2)

	statementList := ""
	for _, statement := range b.Statements {
		statementList += statement.String() + "\n"
	}

	expressionList := ""
	for _, expression := range b.Expressions {
		expr := expression.String()
		expr = strings.ReplaceAll(expr, "\n", "\n"+contentIndent)
		expressionList += contentIndent + expr + "\n"
	}

	scopeStr := ""
	if b.IsNamedScope() {
		scopeStr = bodyIndent + "scope: " + b.Scope + "\n"
	}

	loopStr := ""
	if b.IsLoop {
		loopStr = bodyIndent + "loop: true\n"
	}

	return indent + "block (\n" +
		bodyIndent + "id: " + fmt.Sprint(b.Id) + "\n" +
		scopeStr +
		loopStr +
		bodyIndent + "statements: [\n" +
		statementList +
		bodyIndent + "]\n" +
		bodyIndent + "expressions: [\n" +
		expressionList +
		bodyIndent + "]\n" +
		indent + ")"
}

type ConditionalStatement struct {
	NestLevel int
	Inverted  bool
	Id        int
	WhenTrue  *Block
	WhenFalse *Block
}

func (ConditionalStatement) isStatement() {}

func (c ConditionalStatement) String() string {
	indent := strings.Repeat("\t", c.NestLevel*2)
	bodyIndent := strings.Repeat("\t", c.NestLevel*2+1)

	whenTrueStr := ""
	if c.WhenTrue != nil {
		whenTrueStr = c.WhenTrue.String()
	}

	whenFalseStr := ""
	if c.WhenFalse != nil {
		whenFalseStr = c.WhenFalse.String()
	}

	invertedStr := ""
	if c.Inverted {
		invertedStr = bodyIndent + "inverted: true\n"
	}

	return indent + "conditional (\n" +
		bodyIndent + "id: " + fmt.Sprint(c.Id) + "\n" +
		invertedStr +
		bodyIndent + "when-true:\n" +
		whenTrueStr + "\n" +
		bodyIndent + "when-false:\n" +
		whenFalseStr + "\n" +
		indent + ")"
}

type FunctionStatement struct {
	Name    string
	Args    int
	Returns int
	Body    *Block
}

func (FunctionStatement) isStatement() {}

func (f FunctionStatement) String() string {
	indent := strings.Repeat("\t", 0)
	bodyIndent := strings.Repeat("\t", 1)

	return indent + "function (\n" +
		bodyIndent + "name: " + f.Name + "\n" +
		bodyIndent + "args: " + fmt.Sprint(f.Args) + "\n" +
		bodyIndent + "returns: " + fmt.Sprint(f.Returns) + "\n" +
		bodyIndent + "body:\n" +
		f.Body.String() + "\n" +
		indent + ")"
}

type ExecuteBlockExpression struct {
	BlockId int
}

func (ExecuteBlockExpression) isExpression() {}

func (e ExecuteBlockExpression) String() string {
	return "execute-block " + fmt.Sprint(e.BlockId)
}

type ExecuteConditionalExpression struct {
	ConditionalId int
}

func (ExecuteConditionalExpression) isExpression() {}

func (e ExecuteConditionalExpression) String() string {
	return "execute-conditional " + fmt.Sprint(e.ConditionalId)
}
