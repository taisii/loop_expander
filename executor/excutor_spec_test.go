package executor_test

import (
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

func TestExecute(t *testing.T) {
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
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "add", Operands: []interface{}{"x", 0}}},
							{PC: 1, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"y"}}, Value: 2},
						},
						PathCond: executor.SymbolicExpr{Op: "==", Operands: []interface{}{"x", 0}},
					},
				},
			},
			ExpectError: false,
		},
		// 投機的実行（コミット成功）
		{
			Name: "Speculative branch prediction success",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"r1", "10"}},
				{Mnemonic: "add", Operands: []string{"r2", "r2", "1"}},
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
					PC: 10,
					Registers: map[string]interface{}{
						"r1": 0,
						"r2": 1,
					},
					Memory: map[int]interface{}{},
					Trace:  executor.Trace{},
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
			finalConfigs, err := executor.SpecExecute(testCase.Program, testCase.InitialConfig, testCase.MaxSteps)

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
					}
				}
			}
		})
	}
}
