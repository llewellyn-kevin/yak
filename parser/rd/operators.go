package rd

import (
	"llewellyn-kevin/yak/ast"
)

func (p *RdParser) parseAdd() ast.Expression {
	return ast.AddExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseSubtract() ast.Expression {
	return ast.SubtractExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseMultiply() ast.Expression {
	return ast.MultiplyExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseDivide() ast.Expression {
	return ast.DivideExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseModulo() ast.Expression {
	return ast.ModuloExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parsePower() ast.Expression {
	return ast.PowerExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseLessThan() ast.Expression {
	return ast.LessThanExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseGreaterThan() ast.Expression {
	return ast.GreaterThanExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseLessThanEqualTo() ast.Expression {
	return ast.LessThanEqualToExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseGreaterThanEqualTo() ast.Expression {
	return ast.GreaterThanEqualToExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseEqualTo() ast.Expression {
	return ast.EqualToExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseBor() ast.Expression {
	return ast.BorExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseBand() ast.Expression {
	return ast.BandExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseXor() ast.Expression {
	return ast.XorExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseIncrement() ast.Expression {
	return ast.IncrementExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseDecrement() ast.Expression {
	return ast.DecrementExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseLeftShift() ast.Expression {
	return ast.LeftShiftExpression{
		Token: p.currentToken,
	}
}

func (p *RdParser) parseRightShift() ast.Expression {
	return ast.RightShiftExpression{
		Token: p.currentToken,
	}
}
