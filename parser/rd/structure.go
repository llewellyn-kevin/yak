package rd

import (
	"fmt"
	"llewellyn-kevin/yak/lexer"
	"strings"
)

type Program struct {
	MainBlock *Block
}

func (p Program) String() string {
	return p.MainBlock.String()
}

type Statement interface {
	isStatement()
	String() string
}

type Expression interface {
	isExpression()
	String() string
}

type Block struct {
	id          int
	nestLevel   int
	Scope       string
	IsLoop      bool
	Statements  []Statement
	Expressions []Expression
}

func (Block) isStatement() {}

func (b Block) isRoot() bool {
	return b.id == 0
}

func (b Block) isNamedScope() bool {
	return b.Scope != ""
}

func (p *RdParser) parseBlock(nestLevel int) *Block {
	block := &Block{
		id:        p.blockId,
		nestLevel: nestLevel,
	}
	p.blockId++

	if !block.isRoot() {
		// Check if block has a named scope
		if p.expectPeek(lexer.IDENT) {
			p.nextToken()
			// Check if named scope block is a for loop
			if p.expectPeek(lexer.FOR) {
				block.Scope = p.currentToken.Literal
				block.IsLoop = true
				p.nextToken()
				if p.expectPeek(lexer.COLON) {
					p.nextToken()
				}
			}

			// Not a for loop, so check if this is a named scope, or just an
			// identifier expression
			if p.expectPeek(lexer.COLON) {
				block.Scope = p.currentToken.Literal
				p.nextToken()
			}
		}

		// Check if the unnamed scope is a for loop
		if p.expectPeek(lexer.FOR) {
			p.nextToken()
			block.IsLoop = true
			if p.expectPeek(lexer.COLON) {
				p.nextToken()
			}
		}

		if !p.expectCurrent(lexer.IDENT) {
			p.nextToken()
		}
	}

	for {
		switch {
		case p.expectCurrent(lexer.LBRACE):
			nestedBlock := p.parseBlock(nestLevel + 1)
			block.Statements = append(block.Statements, nestedBlock)
			block.Expressions = append(block.Expressions, ExecuteBlockExpression{
				BlockId: nestedBlock.id,
			})
		case p.expectCurrent(lexer.RBRACE):
			return block
		case p.expectCurrent(lexer.EOF):
			// TODO: If not root, this is an error, but for now just return the block
			return block
		default:
			if expression := p.parseExpression(); expression != nil {
				// TODO: Skipping unkown expressions for now, but we should probably handle this better
				block.Expressions = append(block.Expressions, expression)
			}
		}

		p.nextToken()
	}
}

func (b Block) String() string {
	indent := strings.Repeat("\t", b.nestLevel*2)
	bodyIndent := strings.Repeat("\t", b.nestLevel*2+1)
	contentIndent := strings.Repeat("\t", b.nestLevel*2+2)

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
	if b.isNamedScope() {
		scopeStr = bodyIndent + "scope: " + b.Scope + "\n"
	}

	loopStr := ""
	if b.IsLoop {
		loopStr = bodyIndent + "loop: true\n"
	}

	return indent + "block (\n" +
		bodyIndent + "id: " + fmt.Sprint(b.id) + "\n" +
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

type ExecuteBlockExpression struct {
	BlockId int
}

func (ExecuteBlockExpression) isExpression() {}

func (e ExecuteBlockExpression) String() string {
	return "execute-block " + fmt.Sprint(e.BlockId)
}
