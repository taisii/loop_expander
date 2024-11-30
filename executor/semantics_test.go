package executor_test

import (
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

func TestStep(t *testing.T) {
	testCases := []struct {
		Name            string
		InitialConf     executor.Configuration
		Instruction     assembler.OpCode
		ExpectedConfigs []executor.Configuration
		ExpectError     bool
	}{
		{
			Name: "Symbolic addition (all concreat)",
			InitialConf: executor.Configuration{
				PC: 1,
				Registers: map[string]interface{}{
					"x": executor.SymbolicExpr{
						Op:       "+",
						Operands: []interface{}{5, "y"},
					},
					"y": 3,
				},
			},
			Instruction: assembler.OpCode{
				Mnemonic: "add",
				Operands: []string{"z", "x", "y"}, // z = (5 + y) + y
			},
			ExpectedConfigs: []executor.Configuration{
				{
					PC: 2,
					Registers: map[string]interface{}{
						"x": executor.SymbolicExpr{
							Op:       "+",
							Operands: []interface{}{5, "y"},
						},
						"y": 3,
						"z": 11,
					},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								PC:   1,
								Type: executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{
									Op:       "var",
									Operands: []interface{}{"z"},
								},
								Value: 11,
							},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "Conditional branch with symbolic register",
			InitialConf: executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"x": executor.SymbolicExpr{
						Op:       ">",
						Operands: []interface{}{10, "y"},
					},
					"y": 5,
				},
				Memory: map[int]interface{}{},
			},
			Instruction: assembler.OpCode{
				Mnemonic: "beqz",
				Operands: []string{"x", "100"}, // 10 > y evaluates symbolically
			},
			ExpectedConfigs: []executor.Configuration{
				{
					PC: 1,
					Registers: map[string]interface{}{
						"x": executor.SymbolicExpr{
							Op:       ">",
							Operands: []interface{}{10, "y"},
						},
						"y": 5,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								PC:   0,
								Type: executor.ObsTypePC,
								Value: executor.SymbolicExpr{
									Op:       "!=",
									Operands: []interface{}{1, 0},
								},
							},
						},
						PathCond: executor.SymbolicExpr{
							Op:       "!=",
							Operands: []interface{}{1, 0},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "Conditional branch with symbolic register - Unresolved condition",
			InitialConf: executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"y": 5,
				},
				Memory: map[int]interface{}{},
			},
			Instruction: assembler.OpCode{
				Mnemonic: "beqz",
				Operands: []string{"x", "100"}, // Unable to resolve
			},
			ExpectedConfigs: []executor.Configuration{
				// Case 1: Condition is true (x == 0)
				{
					PC: 100,
					Registers: map[string]interface{}{
						"y": 5,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								PC:   0,
								Type: executor.ObsTypePC,
								Value: executor.SymbolicExpr{
									Op: "==",
									Operands: []interface{}{executor.SymbolicExpr{
										Op:       "symbol",
										Operands: []interface{}{"x"},
									}, 0},
								},
							},
						},
						PathCond: executor.SymbolicExpr{
							Op: "==",
							Operands: []interface{}{executor.SymbolicExpr{
								Op:       "symbol",
								Operands: []interface{}{"x"},
							}, 0},
						},
					},
				},
				// Case 2: Condition is false (x != 0)
				{
					PC: 1,
					Registers: map[string]interface{}{
						"y": 5,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								PC:   0,
								Type: executor.ObsTypePC,
								Value: executor.SymbolicExpr{
									Op: "!=",
									Operands: []interface{}{executor.SymbolicExpr{
										Op:       "symbol",
										Operands: []interface{}{"x"},
									}, 0},
								},
							},
						},
						PathCond: executor.SymbolicExpr{
							Op: "!=",
							Operands: []interface{}{executor.SymbolicExpr{
								Op:       "symbol",
								Operands: []interface{}{"x"},
							}, 0},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "Load operations with symbolic expressions",
			InitialConf: executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"x": 2, // レジスタ x にインデックス値を設定
				},
				Memory: map[int]interface{}{
					10: 100, // メモリアドレス 10 に値 100 を格納 (配列の一部と考える)
					14: 200, // メモリアドレス 14 に値 200 を格納
				},
			},
			Instruction: assembler.OpCode{
				Mnemonic: "load",
				Operands: []string{"y", "10+x*2"}, // メモリアドレスを計算して読み取り
			},
			ExpectedConfigs: []executor.Configuration{
				// Expected configuration after executing the load instruction
				{
					PC: 1,
					Registers: map[string]interface{}{
						"x": 2,
						"y": 200, // メモリアドレス 10+(2*2)=14 の値が y に格納される
					},
					Memory: map[int]interface{}{
						10: 100,
						14: 200,
					},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								PC:      0,
								Type:    executor.ObsTypeLoad,
								Address: 14,
								Value:   200, // 読み取られた値
							},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "Store operation with symbolic index",
			InitialConf: executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"x": 1,   // レジスタ x にインデックス値を設定
					"z": 300, // レジスタ z に格納する値
				},
				Memory: map[int]interface{}{
					10: 100, // メモリアドレス 10 に値 100 を格納
					12: 200, // メモリアドレス 12 に値 200 を格納
				},
			},
			Instruction: assembler.OpCode{
				Mnemonic: "store",
				Operands: []string{"z", "10+x*2"}, // メモリアドレスを計算して z の値を格納
			},
			ExpectedConfigs: []executor.Configuration{
				// Expected configuration after executing the store instruction
				{
					PC: 1,
					Registers: map[string]interface{}{
						"x": 1,
						"z": 300, // レジスタ z の値は変更なし
					},
					Memory: map[int]interface{}{
						10: 100, // メモリアドレス 10 の値は変更なし
						12: 300, // メモリアドレス 10+(1*2)=12 に z の値 (300) が格納される
					},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								PC:      0,
								Type:    executor.ObsTypeStore,
								Address: 12,
								Value:   300, // 書き込まれた値
							},
						},
					},
				},
			},
			ExpectError: false,
		},

		{
			Name: "Symbolic add",
			InitialConf: executor.Configuration{
				PC:        0,
				Registers: map[string]interface{}{},
				Memory:    map[int]interface{}{},
			},
			Instruction: assembler.OpCode{
				Mnemonic: "add",
				Operands: []string{"z", "x", "y"}, // add z x y
			},
			ExpectedConfigs: []executor.Configuration{
				{
					PC: 1,
					Registers: map[string]interface{}{
						"z": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{
								executor.SymbolicExpr{
									Op:       "symbol",
									Operands: []interface{}{"x"},
								},
								executor.SymbolicExpr{
									Op:       "symbol",
									Operands: []interface{}{"y"},
								},
							},
						},
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								PC:   0,
								Type: executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{
									Op:       "var",
									Operands: []interface{}{"z"},
								},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{
											Op:       "symbol",
											Operands: []interface{}{"x"},
										},
										executor.SymbolicExpr{
											Op:       "symbol",
											Operands: []interface{}{"y"},
										},
									},
								},
							},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "Symbolic jump",
			InitialConf: executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"x": executor.SymbolicExpr{
						Op:       "*",
						Operands: []interface{}{2, "y"},
					},
					"y": 4,
				},
				Memory: map[int]interface{}{},
			},
			Instruction: assembler.OpCode{
				Mnemonic: "jmp",
				Operands: []string{"x"}, // Jump to PC = 2 * y
			},
			ExpectedConfigs: []executor.Configuration{
				{
					PC: 8,
					Registers: map[string]interface{}{
						"x": executor.SymbolicExpr{
							Op:       "*",
							Operands: []interface{}{2, "y"},
						},
						"y": 4,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								PC:   0,
								Type: executor.ObsTypePC,
								Value: executor.SymbolicExpr{
									Op:       "jmp",
									Operands: []interface{}{8},
								},
							},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "Unsupported instruction error",
			InitialConf: executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"x": 1,
				},
				Memory: map[int]interface{}{},
			},
			Instruction: assembler.OpCode{
				Mnemonic: "unknown",
				Operands: []string{"x", "y"},
			},
			ExpectedConfigs: []executor.Configuration{
				{
					PC:        1,
					Registers: map[string]interface{}{},
					Memory:    map[int]interface{}{},
				},
			},
			ExpectError: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			finalConfigs, err := executor.Step(testCase.Instruction, &testCase.InitialConf)

			if testCase.ExpectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}

				for i, finalConfig := range finalConfigs {
					if !executor.CompareConfiguration(testCase.ExpectedConfigs[i], *finalConfig) {
						difference := executor.FormatConfigDifferences(testCase.ExpectedConfigs[i], *finalConfig)
						t.Errorf("Test case '%s' failed: configuration %d did not match expected configuration.\n%s",
							testCase.Name, i+1, difference)
						executor.PrintConfiguration(*finalConfigs[i])
					}
				}
			}
		})
	}
}
