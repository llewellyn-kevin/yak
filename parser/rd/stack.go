package rd

import (
	"llewellyn-kevin/yak/ast"
)

func (p *RdParser) parseSwap() ast.Expression {
	return ast.SwapExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseDup() ast.Expression {
	return ast.DupExpression{
		Token: p.currentToken,
	}
}
