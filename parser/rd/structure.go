package rd

import (
	"llewellyn-kevin/yak/ast"
	"llewellyn-kevin/yak/lexer"
)

func (p *RdParser) parseBlock(prgm *ast.Program, nestLevel int) *ast.Block {
	block := &ast.Block{
		Program:   prgm,
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
			conditional := p.parseConditional(prgm, nestLevel+1, p.expectCurrent(lexer.NOT))
			block.Statements = append(block.Statements, conditional)
			block.Expressions = append(block.Expressions, ast.ExecuteConditionalExpression{
				ConditionalId: conditional.Id,
			})
		case p.expectCurrent(lexer.LBRACE):
			nestedBlock := p.parseBlock(prgm, nestLevel+1)
			block.Statements = append(block.Statements, nestedBlock)
			block.Expressions = append(block.Expressions, ast.ExecuteBlockExpression{
				BlockId: nestedBlock.Id,
			})
		case p.expectCurrent(lexer.RBRACE):
			return block
		case p.expectCurrent(lexer.EOF):
			// TODO: If not root, this is an error, but for now just return the block
			return block
		case p.expectCurrent(lexer.INT) && p.expectPeek(lexer.HASH):
			if !block.IsRoot() {
				panic("Function definitions can only be defined in the root block.")
			}
			intLiteral, ok := p.parseInt().(ast.IntLiteral)
			if !ok {
				panic("expected int literal when trying to parse function")
			}
			p.nextToken()
			fn := p.parseSimpleFunc(prgm, intLiteral.Value)
			prgm.AddFunction(fn)
		default:
			if expression := p.parseExpression(); expression != nil {
				// TODO: Skipping unkown expressions for now, but we should probably handle this better
				block.Expressions = append(block.Expressions, expression)
			}
		}

		p.nextToken()
	}
}

func (p *RdParser) parseSimpleFunc(prgm *ast.Program, argCount int) *ast.FunctionStatement {
	f := &ast.FunctionStatement{
		Args: argCount,
	}
	if !p.expectCurrent(lexer.HASH) {
		panic("Expected # after function argument count.")
	}

	p.nextToken()
	if !p.expectCurrent(lexer.IDENT) {
		panic("Expected function name after #.")
	}
	i, ok := p.parseIdent().(ast.IdentifierExpression)
	if !ok {
		panic("Expected function name to be an identifier expression.")
	}
	f.Name = i.Value

	p.nextToken()
	if !p.expectCurrent(lexer.HASH) {
		panic("Expected # after function name.")
	}

	p.nextToken()
	if !p.expectCurrent(lexer.INT) {
		panic("Expected return count after function name.")
	}
	returnCount, ok := p.parseInt().(ast.IntLiteral)
	if !ok {
		panic("Expected return count to be an integer literal.")
	}
	f.Returns = returnCount.Value

	p.nextToken()
	if !p.expectCurrent(lexer.LBRACE) {
		panic("Expected { after function return count.")
	}

	f.Body = p.parseBlock(prgm, 1)

	return f
}

func (p *RdParser) parseConditional(prgm *ast.Program, nestLevel int, inverted bool) *ast.ConditionalStatement {
	c := &ast.ConditionalStatement{}
	c.NestLevel = nestLevel
	c.Inverted = inverted
	c.Id = p.conditionalId
	p.conditionalId++

	if !p.expectPeek(lexer.LBRACE) {
		// TODO: Add this to error handler
		panic("Expected { after if statement.")
	}

	p.nextToken()
	c.WhenTrue = p.parseBlock(prgm, nestLevel+1)

	if !p.expectPeek(lexer.ELSE) {
		c.WhenFalse = &ast.Block{
			Id:          p.blockId,
			NestLevel:   nestLevel + 1,
			Statements:  []ast.Statement{},
			Expressions: []ast.Expression{},
		}
		p.blockId++
		return c
	}

	p.nextToken()

	if !p.expectPeek(lexer.LBRACE) {
		// TODO: Add this to error handler
		panic("Expected { after else statement")
	}

	p.nextToken()
	c.WhenFalse = p.parseBlock(prgm, nestLevel+1)
	return c
}
