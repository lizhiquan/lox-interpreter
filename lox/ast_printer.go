package lox

import (
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr Expr) string {
	return expr.accept(p).(string)
}

var _ exprVisitor = (*AstPrinter)(nil)

func (p *AstPrinter) visitBinaryExpr(expr *BinaryExpr) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p *AstPrinter) visitGroupingExpr(expr *GroupingExpr) any {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) visitLiteralExpr(expr *LiteralExpr) any {
	if expr.Value.Value == nil {
		return "nil"
	}
	return expr.Value.String()
}

func (p *AstPrinter) visitUnaryExpr(expr *UnaryExpr) any {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) any {
	var builder strings.Builder
	builder.WriteString("(")
	builder.WriteString(name)
	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.accept(p).(string))
	}
	builder.WriteString(")")
	return builder.String()
}
