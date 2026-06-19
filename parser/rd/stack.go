package rd

import "llewellyn-kevin/yak/lexer"

type SwapExpression struct {
	Token lexer.Token
}

func (SwapExpression) isExpression() {}

func (s SwapExpression) String() string {
	return "<>"
}

func (p *RdParser) parseSwap() Expression {
	return SwapExpression{
		Token: p.currentToken,
	}
}

type DupExpression struct {
    Token lexer.Token
}

func (DupExpression) isExpression() {}

func (d DupExpression) String() string {
    return "."
}

func (p *RdParser) parseDup() Expression {
    return DupExpression{
        Token: p.currentToken,
    }
}
