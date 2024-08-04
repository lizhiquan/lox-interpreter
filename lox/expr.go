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

func NewBinaryExpr(left Expr, operator Token, right Expr) *BinaryExpr {
	return &BinaryExpr{Left: left, Operator: operator, Right: right}
}

func (expr *BinaryExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitBinaryExpr(expr)
}

type GroupingExpr struct {
	Expression Expr
}

func NewGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{Expression: expression}
}

func (expr *GroupingExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitGroupingExpr(expr)
}

type LiteralExpr struct {
	Value Literal
}

func NewLiteralExpr(value Literal) *LiteralExpr {
	return &LiteralExpr{Value: value}
}

func (expr *LiteralExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitLiteralExpr(expr)
}

type UnaryExpr struct {
	Operator Token
	Right    Expr
}

func NewUnaryExpr(operator Token, right Expr) *UnaryExpr {
	return &UnaryExpr{Operator: operator, Right: right}
}

func (expr *UnaryExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitUnaryExpr(expr)
}

type VariableExpr struct {
	Name Token
}

func NewVariableExpr(name Token) *VariableExpr {
	return &VariableExpr{Name: name}
}

func (expr *VariableExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitVariableExpr(expr)
}

type AssignExpr struct {
	Name  Token
	Value Expr
}

func NewAssignExpr(name Token, value Expr) *AssignExpr {
	return &AssignExpr{Name: name, Value: value}
}

func (expr *AssignExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitAssignExpr(expr)
}

type LogicalExpr struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func NewLogicalExpr(left Expr, operator Token, right Expr) *LogicalExpr {
	return &LogicalExpr{Left: left, Operator: operator, Right: right}
}

func (expr *LogicalExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitLogicalExpr(expr)
}

type CallExpr struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func NewCallExpr(callee Expr, paren Token, arguments []Expr) *CallExpr {
	return &CallExpr{Callee: callee, Paren: paren, Arguments: arguments}
}

func (expr *CallExpr) accept(visitor exprVisitor) (any, error) {
	return visitor.visitCallExpr(expr)
}

type exprVisitor interface {
	visitBinaryExpr(expr *BinaryExpr) (any, error)
	visitGroupingExpr(expr *GroupingExpr) (any, error)
	visitLiteralExpr(expr *LiteralExpr) (any, error)
	visitUnaryExpr(expr *UnaryExpr) (any, error)
	visitVariableExpr(expr *VariableExpr) (any, error)
	visitAssignExpr(expr *AssignExpr) (any, error)
	visitLogicalExpr(expr *LogicalExpr) (any, error)
	visitCallExpr(expr *CallExpr) (any, error)
}
