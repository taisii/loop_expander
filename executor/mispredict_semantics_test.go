package executor_test

import (
	"fmt"
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
	"github.com/taisii/go-project/utils"
)

func TestAlwaysMispredictStep(t *testing.T) {
	// テストケースの定義
	testCases := []struct {
		Name          string
		Instruction   assembler.OpCode
		InitialConfig executor.Configuration
		MaxSpecDepth  int
		ExpectedConfs []executor.Configuration
		ExpectError   bool
	}{
		{
			Name: "beqz with symbolic condition",
			Instruction: assembler.OpCode{
				Mnemonic: "beqz",
				Operands: []string{"x", "2"}, // 条件が成立すればPC=2にジャンプ
			},
			InitialConfig: executor.Configuration{},
			MaxSpecDepth:  2,
			ExpectedConfs: []executor.Configuration{ // 2つの状態が生成される
				{
					PC: 1, // False branch (条件が成立しない場合に次の命令へ)
					Registers: map[string]interface{}{
						"x": executor.SymbolicExpr{Op: "symbolic", Operands: []interface{}{"x"}},
					},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								Type: executor.ObsTypePC,
								PC:   1,
								Value: executor.SymbolicExpr{
									Op:       "==",
									Operands: []interface{}{"x", 0},
								},
							},
						},
						PathCond: executor.SymbolicExpr{
							Op:       "==",
							Operands: []interface{}{"x", 0},
						},
					},
				},
				{
					PC: 2, // False branch (条件が成立しない場合に次の命令へ)
					Registers: map[string]interface{}{
						"x": executor.SymbolicExpr{Op: "symbolic", Operands: []interface{}{"x"}},
					},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								Type: executor.ObsTypePC,
								PC:   1,
								Value: executor.SymbolicExpr{
									Op:       "!=",
									Operands: []interface{}{"x", 0},
								},
							},
						},
						PathCond: executor.SymbolicExpr{
							Op:       "!=",
							Operands: []interface{}{"x", 0},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "jmp instruction",
			Instruction: assembler.OpCode{
				Mnemonic: "jmp",
				Operands: []string{"3"},
			},
			InitialConfig: executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"x": 1,
				},
				Trace: executor.Trace{},
			},
			MaxSpecDepth: 0,
			ExpectedConfs: []executor.Configuration{ // ジャンプ後の状態
				{
					PC: 3, // PC がジャンプ先に移動
					Registers: map[string]interface{}{
						"x": 1,
					},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								Type: executor.ObsTypePC,
								PC:   3,
							},
						},
					},
				},
			},
			ExpectError: false,
		},
	}

	// テストケースの実行
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			execState := &executor.ExecutionState{
				CurrentConf: testCase.InitialConfig,
				Speculative: []executor.SpeculativeState{},
				Counter:     0,
			}

			newConfs, _, err := executor.AlwaysMispredictStep(testCase.Instruction, execState, testCase.MaxSpecDepth)

			// エラーの確認
			if testCase.ExpectError {
				if err == nil {
					t.Errorf("Test case '%s' failed: expected an error but got none", testCase.Name)
				}
				return
			} else {
				if err != nil {
					t.Errorf("Test case '%s' failed: did not expect an error but got one: %v", testCase.Name, err)
					return
				}
			}

			// 結果の構成を比較
			if len(newConfs) != len(testCase.ExpectedConfs) {
				t.Errorf("Test case '%s' failed: expected %d configurations, but got %d",
					testCase.Name, len(testCase.ExpectedConfs), len(newConfs))
			}

			for i, expectedConf := range testCase.ExpectedConfs {
				actualConf := *newConfs[i]

				if executor.CompareConfiguration(expectedConf, actualConf) {
					differences := utils.FormatConfigDifferences(expectedConf, actualConf)
					t.Errorf("Test case '%s' failed: Trace %d did not match expected trace.\n%s",
						testCase.Name, i+1, differences)

					// 失敗した場合は初期状態と最後の状態を出力
					fmt.Println("=== Debug Output ===")
					utils.PrintTest(testCase.InitialConfig, *newConfs[i])
				}
			}
		})
	}
}
