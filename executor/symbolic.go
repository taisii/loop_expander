package executor

import (
	"fmt"
	"strconv"
)

// Configuration structure
type Configuration struct {
	PC        int                    // Program Counter
	Registers map[string]interface{} // General-purpose registers (can hold symbolic or concrete values)
	Memory    map[int]interface{}    // Memory (address to value, symbolic or concrete)
}

// SymbolicExpr represents a symbolic expression.
type SymbolicExpr struct {
	Op       string        // Operator ("+", "-", ">", etc.)
	Operands []interface{} // Operands (can be integers, strings, or nested SymbolicExpr)
}

// NewConfiguration creates a new Configuration
func NewConfiguration(memory map[int]interface{}, registers map[string]interface{}) *Configuration {
	return &Configuration{
		Memory:    memory,
		Registers: registers,
		PC:        0,
	}
}

// Instruction Evaluation
func evalExpr(expr interface{}, conf *Configuration) (interface{}, error) {
	switch expression := expr.(type) {
	case int:
		return expression, nil
	case string: // Could be a register or an integer in string form
		if value, ok := conf.Registers[expression]; ok {
			registerValue, err := evalExpr(value, conf)
			if err != nil {
				return value, nil // Return as symbolic if any operand cannot be evaluated
			}
			return registerValue, nil
		}
		if intValue, err := strconv.Atoi(expression); err == nil {
			return intValue, nil
		}
		return nil, fmt.Errorf("unknown symbol: %s", expression)
	case SymbolicExpr:
		// Attempt to evaluate SymbolicExpr if possible
		evaluatedOperands := make([]interface{}, len(expression.Operands))
		for i, operand := range expression.Operands {
			evalOperand, err := evalExpr(operand, conf)
			if err != nil {
				return expression, nil // Return as symbolic if any operand cannot be evaluated
			}
			evaluatedOperands[i] = evalOperand
		}

		// All operands are concrete, attempt to compute the result
		if allConcrete(evaluatedOperands) {
			result, err := computeConcrete(expression.Op, evaluatedOperands)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		// Return as symbolic if evaluation is not fully concrete
		return SymbolicExpr{
			Op:       expression.Op,
			Operands: evaluatedOperands,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expression)
	}
}

// Helper function to check if all elements in the slice are concrete (int)
func allConcrete(operands []interface{}) bool {
	for _, operand := range operands {
		if _, ok := operand.(int); !ok {
			return false
		}
	}
	return true
}

// Helper function to compute the result of a concrete operation
func computeConcrete(op string, operands []interface{}) (int, error) {
	if len(operands) < 2 {
		return 0, fmt.Errorf("invalid number of operands for operator %s", op)
	}

	// Convert operands to int
	intOperands := make([]int, len(operands))
	for i, operand := range operands {
		intOperands[i] = operand.(int)
	}

	// Perform the operation
	switch op {
	case "+":
		return intOperands[0] + intOperands[1], nil
	case "-":
		return intOperands[0] - intOperands[1], nil
	case "*":
		return intOperands[0] * intOperands[1], nil
	case "/":
		if intOperands[1] == 0 {
			return 0, fmt.Errorf("division by zero")
		}
		return intOperands[0] / intOperands[1], nil
	case "<":
		if intOperands[0] < intOperands[1] {
			return 1, nil
		}
		return 0, nil
	case ">":
		if intOperands[0] > intOperands[1] {
			return 1, nil
		}
		return 0, nil
	case "==":
		if intOperands[0] == intOperands[1] {
			return 1, nil
		}
		return 0, nil
	case "!=":
		if intOperands[0] != intOperands[1] {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported operator: %s", op)
	}
}
