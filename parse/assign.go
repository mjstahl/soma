package parse

import (
	"github.com/socialmachines/soma/ast"
	"github.com/socialmachines/soma/scan"
)

// assignment :=
//   IDENT ASSIGN expression
//
func (p *Parser) parseAssignment(first string) *ast.Assign {
	assign := &ast.Assign{}
	assign.Target = first

	p.expect(scan.ASSIGN)

	assign.Expr = p.parseExpr()
	return assign
}
