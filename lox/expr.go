package lox

// expression     → literal
//                | unary
//                | binary
//                | grouping ;

// literal        → NUMBER | STRING | "true" | "false" | "nil" ;
// grouping       → "(" expression ")" ;
// unary          → ( "-" | "!" ) expression ;
// binary         → expression operator expression ;
// operator       → "==" | "!=" | "<" | "<=" | ">" | ">="
//                | "+"  | "-"  | "*" | "/" ;

type Expr interface {
	accept(visitor exprVisitor) any
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (expr *BinaryExpr) accept(visitor exprVisitor) any {
	return visitor.visitBinaryExpr(expr)
}

type GroupingExpr struct {
	Expression Expr
}

func (expr *GroupingExpr) accept(visitor exprVisitor) any {
	return visitor.visitGroupingExpr(expr)
}

type LiteralExpr struct {
	Value Literal
}

func (expr *LiteralExpr) accept(visitor exprVisitor) any {
	return visitor.visitLiteralExpr(expr)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (expr *UnaryExpr) accept(visitor exprVisitor) any {
	return visitor.visitUnaryExpr(expr)
}

type exprVisitor interface {
	visitBinaryExpr(expr *BinaryExpr) any
	visitGroupingExpr(expr *GroupingExpr) any
	visitLiteralExpr(expr *LiteralExpr) any
	visitUnaryExpr(expr *UnaryExpr) any
}
