package rd

import (
	"llewellyn-kevin/yak/ast"
	"llewellyn-kevin/yak/lexer"
)

func (p *RdParser) parseBlock(nestLevel int) *ast.Block {
	block := &ast.Block{
		Id:        p.blockId,
		NestLevel: nestLevel,
	}
	p.blockId++

	if !block.IsRoot() {
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
		case p.expectCurrent(lexer.IF) || p.expectCurrent(lexer.NOT):
			conditional := p.parseConditional(nestLevel+1, p.expectCurrent(lexer.NOT))
			block.Statements = append(block.Statements, conditional)
			block.Expressions = append(block.Expressions, ast.ExecuteConditionalExpression{
				ConditionalId: conditional.Id,
			})
		case p.expectCurrent(lexer.LBRACE):
			nestedBlock := p.parseBlock(nestLevel + 1)
			block.Statements = append(block.Statements, nestedBlock)
			block.Expressions = append(block.Expressions, ast.ExecuteBlockExpression{
				BlockId: nestedBlock.Id,
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

func (p *RdParser) parseConditional(nestLevel int, inverted bool) (c ast.ConditionalStatement) {
	c.NestLevel = nestLevel
	c.Inverted = inverted
	c.Id = p.conditionalId
	p.conditionalId++

	if !p.expectPeek(lexer.LBRACE) {
		// TODO: Add this to error handler
		panic("Expected { after if statement.")
	}

	p.nextToken()
	c.WhenTrue = p.parseBlock(nestLevel + 1)

	if !p.expectPeek(lexer.ELSE) {
		c.WhenFalse = &ast.Block{
			Id:          p.blockId,
			NestLevel:   nestLevel + 1,
			Statements:  []ast.Statement{},
			Expressions: []ast.Expression{},
		}
		p.blockId++
		return
	}

	p.nextToken()

	if !p.expectPeek(lexer.LBRACE) {
		// TODO: Add this to error handler
		panic("Expected { after else statement")
	}

	p.nextToken()
	c.WhenFalse = p.parseBlock(nestLevel + 1)
	return
}
