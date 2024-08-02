package lox

import (
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr Expr) string {
	s, _ := expr.accept(p)
	return s.(string)
}

var _ exprVisitor = (*AstPrinter)(nil)

func (p *AstPrinter) visitBinaryExpr(expr *BinaryExpr) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (p *AstPrinter) visitGroupingExpr(expr *GroupingExpr) (any, error) {
	return p.parenthesize("group", expr.Expression), nil
}

func (p *AstPrinter) visitLiteralExpr(expr *LiteralExpr) (any, error) {
	if expr.Value.Value == nil {
		return "nil", nil
	}
	return expr.Value.String(), nil
}

func (p *AstPrinter) visitUnaryExpr(expr *UnaryExpr) (any, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (p *AstPrinter) visitVariableExpr(expr *VariableExpr) (any, error) {
	return expr.Name.Lexeme, nil
}

func (p *AstPrinter) visitAssignExpr(expr *AssignExpr) (any, error) {
	return p.parenthesize(expr.Name.Lexeme+" =", expr.Value), nil
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) any {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		s, _ := expr.accept(p)
		builder.WriteString(s.(string))
	}
	builder.WriteString(")")
	return builder.String()
}
