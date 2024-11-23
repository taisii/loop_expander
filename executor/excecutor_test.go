package executor_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

// テストケース用構造体
type TestCase struct {
	Name           string                  // テストケースの名前
	Program        []assembler.OpCode      // プログラム
	InitialConfig  *executor.Configuration // 初期状態
	MaxSteps       int                     // 最大ステップ数
	ExpectedTraces []executor.Trace        // 期待されるトレース
	ExpectError    bool                    // エラーを期待するか
}

func TestExecuteProgram(t *testing.T) {
	// テストケースの定義
	testCases := []TestCase{
		{
			Name: "Simple Program with One Branch",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"x", "3"}},
				{Mnemonic: "mov", Operands: []string{"y", "1"}},
				{Mnemonic: "jmp", Operands: []string{"5"}},
				{Mnemonic: "mov", Operands: []string{"y", "2"}},
			},
			InitialConfig: &executor.Configuration{
				Registers: make(map[string]interface{}),
				PC:        0,
				Trace:     executor.Trace{Observations: []executor.Observation{}},
			},
			MaxSteps: 10,
			ExpectedTraces: []executor.Trace{
				// 分岐の真側のトレース
				{
					Observations: []executor.Observation{
						{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
						{PC: 3, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"y"}}, Value: 2},
					},
					PathCond: executor.SymbolicExpr{
						Op:       "&&",
						Operands: []interface{}{executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
					},
				},
				// 分岐の偽側のトレース
				{
					Observations: []executor.Observation{
						{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
						{PC: 1, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"y"}}, Value: 1},
						{PC: 2, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{5}}},
					},
					PathCond: executor.SymbolicExpr{
						Op:       "&&",
						Operands: []interface{}{executor.SymbolicExpr{Op: "!=", Operands: []interface{}{"x", 0}}},
					},
				},
			},

			ExpectError: false,
		},
		{
			Name: "Program Exceeding Max Steps",
			Program: []assembler.OpCode{
				{Mnemonic: "mov", Operands: []string{"x", "0"}},
				{Mnemonic: "mov", Operands: []string{"y", "0"}},
				{Mnemonic: "jmp", Operands: []string{"1"}},
			},
			InitialConfig: &executor.Configuration{
				Registers: make(map[string]interface{}),
				PC:        0,
				Trace:     executor.Trace{Observations: []executor.Observation{}},
			},
			MaxSteps:       5,
			ExpectedTraces: []executor.Trace{},
			ExpectError:    false,
		},
		{
			Name: "Nested If Statements",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"x", "3"}}, // 条件: x == 0
				{Mnemonic: "mov", Operands: []string{"y", "1"}},  // 偽側
				{Mnemonic: "jmp", Operands: []string{"7"}},       // 偽側の終了
				{Mnemonic: "beqz", Operands: []string{"y", "6"}}, // 真側, 条件: y == 0
				{Mnemonic: "mov", Operands: []string{"z", "2"}},  // 真側の中の偽側
				{Mnemonic: "jmp", Operands: []string{"7"}},       // 真側の中の偽側の終了
				{Mnemonic: "mov", Operands: []string{"z", "3"}},  // 真側の中の真側
			},
			InitialConfig: &executor.Configuration{
				Registers: make(map[string]interface{}),
				PC:        0,
				Trace:     executor.Trace{Observations: []executor.Observation{}},
			},
			MaxSteps: 10,
			ExpectedTraces: []executor.Trace{
				// 真側のトレース (x == 0, y == 0)
				{
					Observations: []executor.Observation{
						{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
						{PC: 3, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"y", 0}}},
						{PC: 6, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"z"}}, Value: 3},
					},
					PathCond: executor.SymbolicExpr{
						Op: "&&",
						Operands: []interface{}{
							executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}},
							executor.SymbolicExpr{Op: "==", Operands: []interface{}{"y", 0}},
						},
					},
				},
				// 偽側のトレース (x != 0)
				{
					Observations: []executor.Observation{
						{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
						{PC: 1, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"y"}}, Value: 1},
						{PC: 2, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
					},
					PathCond: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{"x", 0}},
				},
				// 偽側のトレース (x == 0, y != 0)
				{
					Observations: []executor.Observation{
						{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
						{PC: 3, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"y", 0}}},
						{PC: 4, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"z"}}, Value: 2},
						{PC: 5, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
					},
					PathCond: executor.SymbolicExpr{
						Op: "&&",
						Operands: []interface{}{
							executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}},
							executor.SymbolicExpr{Op: "!=", Operands: []interface{}{"y", 0}},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "Infinite Loop in One Branch",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"x", "3"}}, // 条件: x == 0
				{Mnemonic: "mov", Operands: []string{"y", "1"}},  // 偽側
				{Mnemonic: "jmp", Operands: []string{"0"}},       // 無限ループ
				{Mnemonic: "mov", Operands: []string{"z", "2"}},  // 真側
			},
			InitialConfig: &executor.Configuration{
				Registers: map[string]interface{}{"x": 0},
				PC:        0,
				Trace:     executor.Trace{Observations: []executor.Observation{}},
			},
			MaxSteps: 10,
			ExpectedTraces: []executor.Trace{
				// 無限ループが発生しないため真側のみ
				{
					Observations: []executor.Observation{
						{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
						{PC: 3, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"z"}}, Value: 2},
					},
					PathCond: executor.SymbolicExpr{
						Op:       "&&",
						Operands: []interface{}{executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
					},
				},
			},
			ExpectError: false,
		},
	}

	// テストケースの実行
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			finalConfigs, err := executor.ExecuteProgram(testCase.Program, testCase.InitialConfig, testCase.MaxSteps)

			// エラーの確認
			if testCase.ExpectError {
				if err == nil {
					t.Errorf("Test case '%s' failed: expected an error but got none", testCase.Name)
				}
			} else {
				if err != nil {
					t.Errorf("Test case '%s' failed: did not expect an error but got one: %v", testCase.Name, err)
					return
				}

				// 終了状態の数が一致するか確認
				if len(finalConfigs) != len(testCase.ExpectedTraces) {
					t.Errorf("Test case '%s' failed: expected %d final configurations, but got %d",
						testCase.Name, len(testCase.ExpectedTraces), len(finalConfigs))
					// 失敗した場合は初期状態と最後の状態を出力
					fmt.Println("=== Debug Output ===")
					for _, config := range finalConfigs {
						executor.PrettyPrint(*testCase.InitialConfig, *config)
					}
					return
				}

				// トレースの内容が一致するか確認
				for i, expectedTrace := range testCase.ExpectedTraces {
					actualTrace := finalConfigs[i].Trace // Configuration 内の Trace を取得
					if !compareTraces(expectedTrace, actualTrace) {
						differences := CompareAndPrintDifferences(expectedTrace, actualTrace)
						t.Errorf("Test case '%s' failed: Trace %d did not match expected trace.\n%s",
							testCase.Name, i+1, differences)

						// 失敗した場合は初期状態と最後の状態を出力
						fmt.Println("=== Debug Output ===")
						executor.PrettyPrint(*testCase.InitialConfig, *finalConfigs[i])
					}
				}
			}
		})
	}

}

func compareTraces(expected, actual executor.Trace) bool {
	// 観測数が一致しない場合
	if len(expected.Observations) != len(actual.Observations) {
		return false
	}

	// 各観測の比較
	for i, expObs := range expected.Observations {
		actObs := actual.Observations[i]

		if expObs.PC != actObs.PC ||
			expObs.Type != actObs.Type ||
			!executor.CompareSymbolicExpr(expObs.Address, actObs.Address) ||
			!executor.CompareSymbolicExpr(expObs.Value, actObs.Value) {
			return false
		}
	}

	return true
}

func CompareAndPrintDifferences(expected, actual executor.Trace) string {
	var sb strings.Builder
	sb.WriteString("Differences between expected and actual traces:\n")

	// 観測数の違い
	if len(expected.Observations) != len(actual.Observations) {
		sb.WriteString(fmt.Sprintf("- Observation count mismatch: expected %d, got %d\n",
			len(expected.Observations), len(actual.Observations)))
	} else {
		// 各観測の比較
		for i := 0; i < len(expected.Observations); i++ {
			expectedObs := expected.Observations[i]
			actualObs := actual.Observations[i]

			if expectedObs.PC != actualObs.PC {
				sb.WriteString(fmt.Sprintf("- Mismatch at observation %d (PC): expected %d, got %d\n",
					i+1, expectedObs.PC, actualObs.PC))
			}
			if expectedObs.Type != actualObs.Type {
				sb.WriteString(fmt.Sprintf("- Mismatch at observation %d (Type): expected %s, got %s\n",
					i+1, expectedObs.Type, actualObs.Type))
			}
			if !executor.CompareSymbolicExpr(expectedObs.Address, actualObs.Address) {
				sb.WriteString(fmt.Sprintf("- Mismatch at observation %d (Address):\n", i+1))
				sb.WriteString(fmt.Sprintf("  Expected: %+v\n", expectedObs.Address))
				sb.WriteString(fmt.Sprintf("  Actual:   %+v\n", actualObs.Address))
			}
			if !executor.CompareSymbolicExpr(expectedObs.Value, actualObs.Value) {
				sb.WriteString(fmt.Sprintf("- Mismatch at observation %d (Value):\n", i+1))
				sb.WriteString(fmt.Sprintf("  Expected: %+v\n", expectedObs.Value))
				sb.WriteString(fmt.Sprintf("  Actual:   %+v\n", actualObs.Value))
			}
		}
	}

	// パス条件の違い
	if !executor.CompareSymbolicExpr(expected.PathCond, actual.PathCond) {
		sb.WriteString("- Path condition mismatch:\n")
		sb.WriteString(fmt.Sprintf("  Expected: %+v\n", expected.PathCond))
		sb.WriteString(fmt.Sprintf("  Actual:   %+v\n", actual.PathCond))
	}

	return sb.String()
}
