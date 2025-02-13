package lox

import (
	"fmt"
)

type Interpreter struct {
	globals     *Environment
	environment *Environment
	locals      map[Expr]int
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment()
	globals.Define("clock", Clock{})

	return &Interpreter{
		globals:     globals,
		environment: globals,
		locals:      make(map[Expr]int),
	}
}

func (i *Interpreter) Evaluate(expr Expr) (any, error) {
	return expr.accept(i)
}

func (i *Interpreter) Interpret(statements []Stmt) error {
	for _, stmt := range statements {
		if err := i.execute(stmt); err != nil {
			return err
		}
	}

	return nil
}

var _ exprVisitor = (*Interpreter)(nil)

func (i *Interpreter) visitBinaryExpr(expr *BinaryExpr) (any, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case PLUS:
		if left, ok := left.(string); ok {
			if right, ok := right.(string); ok {
				return left + right, nil
			}
		}

		if left, ok := left.(float64); ok {
			if right, ok := right.(float64); ok {
				return left + right, nil
			}
		}

		return nil, NewRuntimeError(expr.Operator, "Operands must be two numbers or two strings.")
	case MINUS:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case STAR:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case SLASH:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case GREATER:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case GREATER_EQUAL:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case LESS:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case LESS_EQUAL:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case BANG_EQUAL:
		return left != right, nil
	case EQUAL_EQUAL:
		return left == right, nil
	}

	return nil, nil
}

func (i *Interpreter) visitGroupingExpr(expr *GroupingExpr) (any, error) {
	return i.Evaluate(expr.Expression)
}

func (i *Interpreter) visitLiteralExpr(expr *LiteralExpr) (any, error) {
	return expr.Value.Value, nil
}

func (i *Interpreter) visitUnaryExpr(expr *UnaryExpr) (any, error) {
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case MINUS:
		if err := i.checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		}
		return -right.(float64), nil
	case BANG:
		return !i.isTruthy(right), nil
	default:
		return nil, nil
	}
}

func (i *Interpreter) visitVariableExpr(expr *VariableExpr) (any, error) {
	return i.lookUpVariable(expr.Name, expr)
}

func (i *Interpreter) lookUpVariable(name Token, expr Expr) (any, error) {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.GetAt(distance, name.Lexeme), nil
	}

	return i.globals.Get(name)
}

func (i *Interpreter) visitAssignExpr(expr *AssignExpr) (any, error) {
	value, err := i.Evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	if distance, ok := i.locals[expr]; ok {
		i.environment.AssignAt(distance, expr.Name, value)
	} else {
		if err := i.globals.Assign(expr.Name, value); err != nil {
			return nil, err
		}
	}

	return value, nil
}

func (i *Interpreter) visitLogicalExpr(expr *LogicalExpr) (any, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.Evaluate(expr.Right)
}

func (i *Interpreter) visitCallExpr(expr *CallExpr) (any, error) {
	callee, err := i.Evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	callable, ok := callee.(Callable)
	if !ok {
		return nil, NewRuntimeError(expr.Paren, "Can only call functions and classes.")
	}

	if len(expr.Arguments) != callable.Arity() {
		return nil, NewRuntimeError(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", callable.Arity(), len(expr.Arguments)))
	}

	arguments := make([]any, len(expr.Arguments))
	for idx, arg := range expr.Arguments {
		val, err := i.Evaluate(arg)
		if err != nil {
			return nil, err
		}

		arguments[idx] = val
	}

	return callable.Call(i, arguments)
}

var _ stmtVisitor = (*Interpreter)(nil)

func (i *Interpreter) visitExprStmt(stmt *ExprStmt) (any, error) {
	return i.Evaluate(stmt.Expression)
}

func (i *Interpreter) visitPrintStmt(stmt *PrintStmt) (any, error) {
	val, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	if val == nil {
		val = "nil"
	}

	fmt.Println(val)
	return nil, nil
}

func (i *Interpreter) visitVarDeclStmt(stmt *VarDeclStmt) (any, error) {
	var value any
	if stmt.Initializer != nil {
		var err error
		value, err = i.Evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) visitBlockStmt(stmt *BlockStmt) (any, error) {
	return nil, i.executeBlock(stmt.Statements, NewEnvironmentWithEnclosing(i.environment))
}

func (i *Interpreter) visitIfStmt(stmt *IfStmt) (any, error) {
	value, err := i.Evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(value) {
		return nil, i.execute(stmt.ThenBranch)
	}

	if stmt.ElseBranch != nil {
		return nil, i.execute(stmt.ElseBranch)
	}

	return nil, nil
}

func (i *Interpreter) visitWhileStmt(stmt *WhileStmt) (any, error) {
	for {
		condition, err := i.Evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}

		if !i.isTruthy(condition) {
			break
		}

		if err := i.execute(stmt.Body); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Interpreter) visitFunctionDeclStmt(stmt *FunctionDeclStmt) (any, error) {
	function := NewFunction(stmt, i.environment)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) visitReturnStmt(stmt *ReturnStmt) (any, error) {
	var value any
	if stmt.Value != nil {
		var err error
		value, err = i.Evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, NewReturnError(value)
}

func (i *Interpreter) execute(stmt Stmt) error {
	_, err := stmt.accept(i)
	return err
}

func (i *Interpreter) resolve(expr Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) executeBlock(statements []Stmt, environment *Environment) error {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()

	i.environment = environment
	for _, stmt := range statements {
		if err := i.execute(stmt); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) isTruthy(value any) bool {
	switch val := value.(type) {
	case nil:
		return false
	case bool:
		return val
	default:
		return true
	}
}

func (i *Interpreter) checkNumberOperand(operator Token, operand any) error {
	if _, ok := operand.(float64); ok {
		return nil
	}
	return NewRuntimeError(operator, "Operand must be a number.")
}

func (i *Interpreter) checkNumberOperands(operator Token, left, right any) error {
	if _, ok := left.(float64); ok {
		if _, ok := right.(float64); ok {
			return nil
		}
	}
	return NewRuntimeError(operator, "Operands must be numbers.")
}
