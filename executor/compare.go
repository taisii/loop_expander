package executor

import (
	"reflect"
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

func CompareTraces(expected, actual Trace) bool {
	// 観測数が一致しない場合
	if len(expected.Observations) != len(actual.Observations) {
		return false
	}

	// 各観測の比較
	for i, expObs := range expected.Observations {
		actObs := actual.Observations[i]

		if expObs.PC != actObs.PC ||
			expObs.Type != actObs.Type ||
			!CompareSymbolicExpr(expObs.Address, actObs.Address) ||
			!CompareSymbolicExpr(expObs.Value, actObs.Value) {
			return false
		}
	}

	return true
}

func CompareConfiguration(expected, actual Configuration) bool {
	// プログラムカウンタ (PC) の比較
	if expected.PC != actual.PC {
		return false
	}

	// ステップカウントの比較
	if expected.StepCount != actual.StepCount {
		return false
	}

	// レジスタの比較
	if !CompareRegisters(expected.Registers, actual.Registers) {
		return false
	}

	// メモリの比較
	if !CompareMemory(expected.Memory, actual.Memory) {
		return false
	}

	// トレースの比較
	if !CompareTraces(expected.Trace, actual.Trace) {
		return false
	}

	return true
}

// レジスタを比較
func CompareRegisters(expected, actual map[string]interface{}) bool {
	if len(expected) != len(actual) {
		return false
	}

	for key, expVal := range expected {
		actVal, exists := actual[key]
		if !exists || !CompareSymbolicExpr(expVal, actVal) {
			return false
		}
	}

	return true
}

// メモリを比較
func CompareMemory(expected, actual map[int]interface{}) bool {
	if len(expected) != len(actual) {
		return false
	}

	for addr, expVal := range expected {
		actVal, exists := actual[addr]
		if !exists || !CompareSymbolicExpr(expVal, actVal) {
			return false
		}
	}

	return true
}
