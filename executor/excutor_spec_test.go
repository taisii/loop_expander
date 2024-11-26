package executor_test

import (
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

func TestSpecExecute(t *testing.T) {
	// テストケースの定義
	testCases := []struct {
		Name            string                   // テストケースの名前
		Program         []assembler.OpCode       // プログラム
		InitialConfig   *executor.Configuration  // 初期状態
		MaxSteps        int                      // 最大ステップ数
		ExpectedConfigs []executor.Configuration // 期待される終了状態
		ExpectError     bool                     // エラーを期待するか
	}{
		// 非投機的実行
		{
			Name: "Simple non-speculative execution",
			Program: []assembler.OpCode{
				{Mnemonic: "add", Operands: []string{"r1", "r2", "1"}},
				{Mnemonic: "jmp", Operands: []string{"2"}},
			},
			InitialConfig: &executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"r1": 0,
					"r2": 0,
				},
				Memory: map[int]interface{}{},
				Trace:  executor.Trace{},
			},
			MaxSteps: 10,
			ExpectedConfigs: []executor.Configuration{
				{
					PC: 2,
					Registers: map[string]interface{}{
						"r1": 1,
						"r2": 0,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r1"}}, Value: 1},
							{PC: 1, Type: executor.ObsTypePC, Address: nil, Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{2}}},
						},
						PathCond: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}},
					},
				},
			},
			ExpectError: false,
		},
		// 投機的実行（ロールバック）
		{
			Name: "Rollback on speculative execution misprediction",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"r1", "3"}}, // 投機的実行を開始
				{Mnemonic: "add", Operands: []string{"r2", "r2", "1"}},
				{Mnemonic: "add", Operands: []string{"r3", "r3", "1"}},
			},
			InitialConfig: &executor.Configuration{},
			MaxSteps:      10,
			ExpectedConfigs: []executor.Configuration{
				{ // r1 = 3のとき
					PC: 3,
					Registers: map[string]interface{}{
						"r2": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{executor.SymbolicExpr{
								Op:       "symbol",
								Operands: []interface{}{"r2"}}, 1}},
						"r3": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{executor.SymbolicExpr{
								Op:       "symbol",
								Operands: []interface{}{"r3"}}, 1}},
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStart, Value: 0},
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{"r1", 0}}},
							{PC: 1, Type: executor.ObsTypeRollback, Value: 0},
							{PC: 1,
								Type: executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{
									Op:       "var",
									Operands: []interface{}{"r2"}},
								Value: executor.SymbolicExpr{Op: "+", Operands: []interface{}{executor.SymbolicExpr{
									Op:       "symbol",
									Operands: []interface{}{"r2"}}, 1}}},
							{PC: 2,
								Type: executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{
									Op:       "var",
									Operands: []interface{}{"r3"}},
								Value: executor.SymbolicExpr{Op: "+", Operands: []interface{}{executor.SymbolicExpr{
									Op:       "symbol",
									Operands: []interface{}{"r3"}}, 1}}},
						},
						PathCond: executor.SymbolicExpr{
							Op: "==",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}},
								0,
							},
						},
					},
				},
				{
					PC: 4, // 正常ルート（ロールバック後の状態）
					Registers: map[string]interface{}{
						"r1": 1,
						"r2": 1,
						"r3": 10,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}}},
						},
					},
				},
			},
			ExpectError: false,
		},
		// 投機的実行（ロールバック）
		{
			Name: "Speculative branch rollback",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"r1", "10"}},
				{Mnemonic: "add", Operands: []string{"r2", "r2", "1"}},
			},
			InitialConfig: &executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"r1": 42, // 分岐条件が成立しない
					"r2": 0,
				},
				Memory: map[int]interface{}{},
				Trace:  executor.Trace{},
			},
			MaxSteps: 10,
			ExpectedConfigs: []executor.Configuration{
				{
					PC: 1,
					Registers: map[string]interface{}{
						"r1": 42,
						"r2": 1,
					},
					Memory: map[int]interface{}{},
					Trace:  executor.Trace{},
				},
			},
			ExpectError: false,
		},
		// 無限ループ防止
		{
			Name: "Infinite loop prevention",
			Program: []assembler.OpCode{
				{Mnemonic: "jmp", Operands: []string{"0"}}, // 自己ループ
			},
			InitialConfig: &executor.Configuration{
				PC:        0,
				Registers: map[string]interface{}{},
				Memory:    map[int]interface{}{},
				Trace:     executor.Trace{},
			},
			MaxSteps:        10,
			ExpectedConfigs: nil,
			ExpectError:     true,
		},
		// 無効な命令
		{
			Name: "Unsupported instruction",
			Program: []assembler.OpCode{
				{Mnemonic: "unsupported", Operands: []string{}},
			},
			InitialConfig: &executor.Configuration{
				PC:        0,
				Registers: map[string]interface{}{},
				Memory:    map[int]interface{}{},
				Trace:     executor.Trace{},
			},
			MaxSteps:        10,
			ExpectedConfigs: nil,
			ExpectError:     true,
		},
	}

	// 各テストケースの実行
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// `execute`関数を呼び出し
			finalConfigs, err := executor.SpecExecute(testCase.Program, testCase.InitialConfig, testCase.MaxSteps, 10)

			// エラーの確認
			if testCase.ExpectError {
				if err == nil {
					t.Errorf("Test case '%s' failed: expected an error but got none", testCase.Name)
				}
			} else {
				if err != nil {
					t.Errorf("Test case '%s' failed: unexpected error: %v", testCase.Name, err)
					return
				}

				// 終了状態の数が期待通りか確認
				if len(finalConfigs) != len(testCase.ExpectedConfigs) {
					t.Errorf("Test case '%s' failed: expected %d final configurations, but got %d",
						testCase.Name, len(testCase.ExpectedConfigs), len(finalConfigs))
					return
				}

				// 各終了状態の比較
				for i, expectedConfig := range testCase.ExpectedConfigs {
					if !executor.CompareConfiguration(expectedConfig, *finalConfigs[i]) {
						difference := executor.FormatConfigDifferences(expectedConfig, *finalConfigs[i])
						t.Errorf("Test case '%s' failed: configuration %d did not match expected configuration.\n%s",
							testCase.Name, i+1, difference)
						executor.PrintConfiguration(*finalConfigs[i])
					}
				}
			}
		})
	}
}
