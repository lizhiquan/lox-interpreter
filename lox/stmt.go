package lox

type Stmt interface {
	accept(visitor stmtVisitor) (any, error)
}

type ExprStmt struct {
	Expression Expr
}

func NewExprStmt(expression Expr) *ExprStmt {
	return &ExprStmt{Expression: expression}
}

func (s *ExprStmt) accept(visitor stmtVisitor) (any, error) {
	return visitor.visitExprStmt(s)
}

type PrintStmt struct {
	Expression Expr
}

func NewPrintStmt(expression Expr) *PrintStmt {
	return &PrintStmt{Expression: expression}
}

func (s *PrintStmt) accept(visitor stmtVisitor) (any, error) {
	return visitor.visitPrintStmt(s)
}

type VarDeclStmt struct {
	Name        Token
	Initializer Expr
}

func NewVarDeclStmt(name Token, initializer Expr) *VarDeclStmt {
	return &VarDeclStmt{Name: name, Initializer: initializer}
}

func (s *VarDeclStmt) accept(visitor stmtVisitor) (any, error) {
	return visitor.visitVarDeclStmt(s)
}

type BlockStmt struct {
	Statements []Stmt
}

func NewBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{Statements: statements}
}

func (s *BlockStmt) accept(visitor stmtVisitor) (any, error) {
	return visitor.visitBlockStmt(s)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) *IfStmt {
	return &IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (s *IfStmt) accept(visitor stmtVisitor) (any, error) {
	return visitor.visitIfStmt(s)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func NewWhileStmt(condition Expr, body Stmt) *WhileStmt {
	return &WhileStmt{Condition: condition, Body: body}
}

func (s *WhileStmt) accept(visitor stmtVisitor) (any, error) {
	return visitor.visitWhileStmt(s)
}

type stmtVisitor interface {
	visitExprStmt(stmt *ExprStmt) (any, error)
	visitPrintStmt(stmt *PrintStmt) (any, error)
	visitVarDeclStmt(stmt *VarDeclStmt) (any, error)
	visitBlockStmt(stmt *BlockStmt) (any, error)
	visitIfStmt(stmt *IfStmt) (any, error)
	visitWhileStmt(stmt *WhileStmt) (any, error)
}
