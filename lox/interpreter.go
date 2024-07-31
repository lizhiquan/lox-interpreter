package lox

import "fmt"

type Interpreter struct{}

func (i *Interpreter) Interpret(expr Expr) (string, error) {
	val, err := i.evaluate(expr)
	if err != nil {
		return "", err
	}

	if val == nil {
		return "nil", nil
	}

	return fmt.Sprint(val), nil
}

var _ exprVisitor = (*Interpreter)(nil)

func (i *Interpreter) visitBinaryExpr(expr *BinaryExpr) (any, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.evaluate(expr.Right)
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
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) visitLiteralExpr(expr *LiteralExpr) (any, error) {
	return expr.Value.Value, nil
}

func (i *Interpreter) visitUnaryExpr(expr *UnaryExpr) (any, error) {
	right, err := i.evaluate(expr.Right)
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

func (i *Interpreter) evaluate(expr Expr) (any, error) {
	return expr.accept(i)
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

type RuntimeError struct {
	token   Token
	message string
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{token: token, message: message}
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] %s", e.token.Line, e.message)
}
