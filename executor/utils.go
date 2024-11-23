package executor

import (
	"fmt"
	"reflect"
	"strings"
)

func CompareSymbolicExpr(expected, actual interface{}) bool {
	// 両方とも nil の場合は一致
	if expected == nil && actual == nil {
		return true
	}
	// どちらか一方が nil の場合は不一致
	if expected == nil || actual == nil {
		return false
	}

	// 型アサーションで `SymbolicExpr` にキャスト可能か確認
	expExpr, okExp := expected.(*SymbolicExpr)
	actExpr, okAct := actual.(*SymbolicExpr)

	// 両方とも `*SymbolicExpr` の場合、再帰的に比較
	if okExp && okAct {
		// 演算子とオペランド数が異なる場合は不一致
		if expExpr.Op != actExpr.Op || len(expExpr.Operands) != len(actExpr.Operands) {
			return false
		}
		// オペランドごとに比較（再帰的）
		for i, expOperand := range expExpr.Operands {
			if !CompareSymbolicExpr(expOperand, actExpr.Operands[i]) {
				return false
			}
		}
		return true
	}

	// 型アサーションが失敗した場合、直接値として比較（文字列や整数など）
	return reflect.DeepEqual(expected, actual)
}

// PrettyPrint トレース全体をフォーマットして出力する
func PrettyPrint(initialConfig, finalConfig Configuration) {
	// Assignments 出力
	fmt.Println("Assignments:")
	printMapStringInterface(finalConfig.Registers, "  ")
	fmt.Println()

	// 初期状態の出力
	fmt.Println("initial conf:")
	printConfiguration(initialConfig)

	// トレースの出力
	fmt.Println("\ntrace:")
	for _, obs := range finalConfig.Trace.Observations {
		printObservation(obs)
	}

	// 最終状態の出力
	fmt.Println("\nfinal conf:")
	printConfiguration(finalConfig)

	// Path Condition の出力
	fmt.Println("\nPath Condition:")
	fmt.Printf("  %s\n", formatSymbolicExpr(finalConfig.Trace.PathCond))
	fmt.Println("===========================")
}

// printConfiguration Configuration を整形して出力
func printConfiguration(config Configuration) {
	fmt.Println("  m=")
	printMapStringInterface(config.Registers, "    ")
	fmt.Println("  a=")
	printMapIntInterface(config.Memory, "    ")
}

// printMapStringInterface map[string]interface{} を整形して出力
func printMapStringInterface(data map[string]interface{}, indent string) {
	for key, value := range data {
		fmt.Printf("%s%s: %s\n", indent, key, formatValue(value))
	}
}

// printMapIntInterface map[int]interface{} を整形して出力
func printMapIntInterface(data map[int]interface{}, indent string) {
	for key, value := range data {
		fmt.Printf("%s%d: %s\n", indent, key, formatValue(value))
	}
}

// printObservation 観測データを整形して出力
func printObservation(obs Observation) {
	fmt.Printf("  PC: %d, Type: %s", obs.PC, obs.Type)
	if obs.Address != nil {
		fmt.Printf(", Address: %s", formatValue(obs.Address))
	}
	if obs.Value != nil {
		fmt.Printf(", Value: %s", formatValue(obs.Value))
	}
	if obs.SpecState != nil {
		fmt.Printf(", SpeculativeState: {ID: %d, Depth: %d, RolledBack: %t}",
			obs.SpecState.ID, obs.SpecState.Depth, obs.SpecState.IsRolledBack)
	}
	fmt.Println()
}

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
