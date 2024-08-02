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

type stmtVisitor interface {
	visitExprStmt(stmt *ExprStmt) (any, error)
	visitPrintStmt(stmt *PrintStmt) (any, error)
	visitVarDeclStmt(stmt *VarDeclStmt) (any, error)
	visitBlockStmt(stmt *BlockStmt) (any, error)
}
