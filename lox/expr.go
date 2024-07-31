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
	accept(visitor exprVisitor) (any, error)
}

type BinaryExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (expr *BinaryExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitBinaryExpr(expr)
}

type GroupingExpr struct {
	Expression Expr
}

func (expr *GroupingExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitGroupingExpr(expr)
}

type LiteralExpr struct {
	Value Literal
}

func (expr *LiteralExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitLiteralExpr(expr)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func (expr *UnaryExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitUnaryExpr(expr)
}

type exprVisitor interface {
	visitBinaryExpr(expr *BinaryExpr) (any, error)
	visitGroupingExpr(expr *GroupingExpr) (any, error)
	visitLiteralExpr(expr *LiteralExpr) (any, error)
	visitUnaryExpr(expr *UnaryExpr) (any, error)
}
