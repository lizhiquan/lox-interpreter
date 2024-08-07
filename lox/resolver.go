package lox

type Resolver struct {
	interpreter     *Interpreter
	scopes          []map[string]bool
	currentFunction FunctionType
}

type FunctionType int

const (
	NONE FunctionType = iota
	FUNCTION
)

func NewResolver(interpreter *Interpreter) *Resolver {
	return &Resolver{
		interpreter:     interpreter,
		scopes:          []map[string]bool{},
		currentFunction: NONE,
	}
}

var (
	_ exprVisitor = (*Resolver)(nil)
	_ stmtVisitor = (*Resolver)(nil)
)

func (r *Resolver) visitBlockStmt(stmt *BlockStmt) (any, error) {
	r.beginScope()

	if err := r.Resolve(stmt.Statements); err != nil {
		return nil, err
	}

	r.endScope()
	return nil, nil
}

func (r *Resolver) visitVarDeclStmt(stmt *VarDeclStmt) (any, error) {
	if err := r.declare(stmt.Name); err != nil {
		return nil, err
	}

	if stmt.Initializer != nil {
		if err := r.resolveExpr(stmt.Initializer); err != nil {
			return nil, err
		}
	}

	r.define(stmt.Name)
	return nil, nil
}

func (r *Resolver) visitVariableExpr(expr *VariableExpr) (any, error) {
	if len(r.scopes) > 0 {
		if ready, ok := r.scopes[len(r.scopes)-1][expr.Name.Lexeme]; ok && !ready {
			return nil, NewRuntimeError(expr.Name, "Can't read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) visitAssignExpr(expr *AssignExpr) (any, error) {
	if err := r.resolveExpr(expr.Value); err != nil {
		return nil, err
	}

	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) visitFunctionDeclStmt(stmt *FunctionDeclStmt) (any, error) {
	if err := r.declare(stmt.Name); err != nil {
		return nil, err
	}

	r.define(stmt.Name)

	return nil, r.resolveFunction(stmt, FUNCTION)
}

func (r *Resolver) visitExprStmt(stmt *ExprStmt) (any, error) {
	return nil, r.resolveExpr(stmt.Expression)
}

func (r *Resolver) visitIfStmt(stmt *IfStmt) (any, error) {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}

	if err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return nil, err
	}

	if stmt.ElseBranch != nil {
		if err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) visitPrintStmt(stmt *PrintStmt) (any, error) {
	return nil, r.resolveExpr(stmt.Expression)
}

func (r *Resolver) visitReturnStmt(stmt *ReturnStmt) (any, error) {
	if r.currentFunction == NONE {
		return nil, NewRuntimeError(stmt.Keyword, "Can't return from top-level code.")
	}

	if stmt.Value != nil {
		if err := r.resolveExpr(stmt.Value); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) visitWhileStmt(stmt *WhileStmt) (any, error) {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return nil, err
	}

	return nil, r.resolveStmt(stmt.Body)
}

func (r *Resolver) visitBinaryExpr(expr *BinaryExpr) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}

	return nil, r.resolveExpr(expr.Right)
}

func (r *Resolver) visitCallExpr(expr *CallExpr) (any, error) {
	if err := r.resolveExpr(expr.Callee); err != nil {
		return nil, err
	}

	for _, argument := range expr.Arguments {
		if err := r.resolveExpr(argument); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *Resolver) visitGroupingExpr(expr *GroupingExpr) (any, error) {
	return nil, r.resolveExpr(expr.Expression)
}

func (r *Resolver) visitLiteralExpr(expr *LiteralExpr) (any, error) {
	return nil, nil
}

func (r *Resolver) visitLogicalExpr(expr *LogicalExpr) (any, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return nil, err
	}

	return nil, r.resolveExpr(expr.Right)
}

func (r *Resolver) visitUnaryExpr(expr *UnaryExpr) (any, error) {
	return nil, r.resolveExpr(expr.Right)
}

func (r *Resolver) Resolve(statements []Stmt) error {
	for _, stmt := range statements {
		if err := r.resolveStmt(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (r *Resolver) resolveStmt(stmt Stmt) error {
	_, err := stmt.accept(r)
	return err
}

func (r *Resolver) resolveExpr(expr Expr) error {
	_, err := expr.accept(r)
	return err
}

func (r *Resolver) resolveFunction(function *FunctionDeclStmt, functionType FunctionType) error {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType
	defer func() { r.currentFunction = enclosingFunction }()

	r.beginScope()

	for _, param := range function.Parameters {
		if err := r.declare(param); err != nil {
			return err
		}

		r.define(param)
	}

	if err := r.Resolve(function.Body); err != nil {
		return err
	}

	r.endScope()
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *Resolver) declare(name Token) error {
	if len(r.scopes) == 0 {
		return nil
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.Lexeme]; ok {
		return NewRuntimeError(name, "Already a variable with this name in this scope.")
	}

	scope[name.Lexeme] = false
	return nil
}

func (r *Resolver) define(name Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.Lexeme]; ok {
			r.interpreter.resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}
