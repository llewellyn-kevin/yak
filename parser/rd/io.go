package rd

import (
	"llewellyn-kevin/yak/ast"
)

func (p *RdParser) parseYakout() ast.YakoutExpression {
	return ast.YakoutExpression{
		Token: p.currentToken,
	}
}
