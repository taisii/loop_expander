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
	Trace     Trace
	StepCount int
}

// SymbolicExpr represents a symbolic expression.
type SymbolicExpr struct {
	Op       string        // Operator ("+", "-", ">", etc.)
	Operands []interface{} // Operands (can be integers, strings, or nested SymbolicExpr)
}

type Trace struct {
	Observations []Observation // 実行過程の観測データのリスト
	PathCond     SymbolicExpr  // このトレースのパス条件（シンボリック形式）
}

type Observation struct {
	PC        int               // プログラムカウンタ（実行中の命令の位置）
	Type      ObsType           // 観測タイプ: load, store, pc, start, rollback, commit
	Address   interface{}       // メモリアクセスの場合のアドレス（シンボリック形式）
	Value     interface{}       // 値の読み取りや書き込みの内容（シンボリック形式）
	SpecState *SpeculativeState // スペキュレーション状態（該当する場合）
}

type ObsType string

const (
	ObsTypeLoad     ObsType = "load"     // メモリ読み取り
	ObsTypeStore    ObsType = "store"    // メモリ書き込み
	ObsTypePC       ObsType = "pc"       // プログラムカウンタの変更
	ObsTypeStart    ObsType = "start"    // スペキュレーション開始
	ObsTypeRollback ObsType = "rollback" // スペキュレーション取り消し
	ObsTypeCommit   ObsType = "commit"   // スペキュレーションのコミット
)

type SpeculativeState struct {
	ID           int  // スペキュレーションの一意の識別子
	Depth        int  // ネストの深さ
	IsRolledBack bool // このスペキュレーションが取り消されたかどうか
}

// NewConfiguration creates a new Configuration
func NewConfiguration(memory map[int]interface{}, registers map[string]interface{}) *Configuration {
	return &Configuration{
		Memory:    memory,
		Registers: registers,
		PC:        0,
	}
}

type SymbolicValue struct {
	Concrete *int          // 具体値がある場合
	Symbolic *SymbolicExpr // シンボリック式がある場合
}

// Instruction Evaluation
func evalExpr(expr interface{}, conf *Configuration) (interface{}, error) {
	switch expression := expr.(type) {
	case int:
		return expression, nil
	case string: // Could be a register or an integer in string form
		// レジスタに存在するか確認
		if value, ok := conf.Registers[expression]; ok {
			// 再帰的に評価
			registerValue, err := evalExpr(value, conf)
			if err != nil {
				// 評価中にエラーがあればシンボリックなまま返す
				return value, nil
			}
			return registerValue, nil
		}

		// 整数値として解析可能か確認
		if intValue, err := strconv.Atoi(expression); err == nil {
			return intValue, nil
		}

		// 未定義の変数はシンボリック変数として扱う
		return SymbolicExpr{
			Op:       "symbol",
			Operands: []interface{}{expression},
		}, nil
	case SymbolicExpr:
		// シンボリック式を評価
		evaluatedOperands := make([]interface{}, len(expression.Operands))
		for i, operand := range expression.Operands {
			evalOperand, err := evalExpr(operand, conf)
			if err != nil {
				// 評価中にエラーがあれば式全体をシンボリックのまま返す
				return expression, nil
			}
			evaluatedOperands[i] = evalOperand
		}

		// オペランドがすべて具体値なら結果を計算
		if allConcrete(evaluatedOperands) {
			result, err := computeConcrete(expression.Op, evaluatedOperands)
			if err != nil {
				return nil, err
			}
			return result, nil
		}

		// オペランドにシンボリックな値が含まれる場合、シンボリック式として返す
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
