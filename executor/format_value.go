package executor

import (
	"fmt"
	"strings"
)

// formatValue 値を適切にフォーマット
func formatValue(value interface{}) string {
	switch v := value.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	case *SymbolicExpr:
		return formatSymbolicExpr(*v)
	case SymbolicExpr:
		return formatSymbolicExpr(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// formatSymbolicExpr シンボリック式を文字列にフォーマット
func formatSymbolicExpr(expr SymbolicExpr) string {
	// 単一のオペランドの場合は括弧で囲まずに出力
	if len(expr.Operands) == 1 {
		return formatValue(expr.Operands[0])
	}

	// 演算を伴う場合はオペランドを演算子で結合し、式全体を () で囲む
	var operands []string
	for _, op := range expr.Operands {
		operands = append(operands, formatValue(op))
	}
	return fmt.Sprintf("(%s)", strings.Join(operands, fmt.Sprintf(" %s ", expr.Op)))
}
