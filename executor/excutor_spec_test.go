package executor_test

import (
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

func TestSpecExecute(t *testing.T) {
	// テストケースの定義
	testCases := []TestCase{
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
		{
			Name: "symbolic single benq ",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"r1", "3"}},
				{Mnemonic: "add", Operands: []string{"r2", "r2", "1"}},
				{Mnemonic: "add", Operands: []string{"r3", "r3", "1"}},
			},
			InitialConfig: &executor.Configuration{},
			MaxSteps:      10,
			ExpectedConfigs: []executor.Configuration{
				{ // r1 == 0のとき
					PC:        3,
					Registers: map[string]interface{}{},
					Memory:    map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStart, Value: 0},
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{
								Op: "symbol", Operands: []interface{}{"r1"}}, 0}}},
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
							{PC: 3, Type: executor.ObsTypeRollback, Value: 0},
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
				{ // r1 != 0のとき
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
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{
								Op: "symbol", Operands: []interface{}{"r1"}}, 0}}},
							{PC: 3, Type: executor.ObsTypeRollback, Value: 0},
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
							Op: "!=",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}},
								0,
							},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "symbolic branch true with 7 add ops",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"r1", "7"}}, // if r1 == 0, jump to PC 8
				{Mnemonic: "add", Operands: []string{"r2", "r2", "1"}},
				{Mnemonic: "add", Operands: []string{"r3", "r3", "2"}},
				{Mnemonic: "add", Operands: []string{"r4", "r4", "3"}},
				{Mnemonic: "add", Operands: []string{"r5", "r5", "4"}},
				{Mnemonic: "add", Operands: []string{"r6", "r6", "5"}},
				{Mnemonic: "add", Operands: []string{"r7", "r7", "6"}},
				{Mnemonic: "add", Operands: []string{"r8", "r8", "7"}},
			},
			InitialConfig: &executor.Configuration{},
			MaxSteps:      100,
			ExpectedConfigs: []executor.Configuration{
				{ // r1 == 0のとき
					PC: 8,
					Registers: map[string]interface{}{
						"r8": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r8"}},
								7,
							},
						},
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStart, Value: 0},
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{
								Op: "!=",
								Operands: []interface{}{
									executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}},
									0,
								},
							}},
							{PC: 1,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r2"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}},
										1,
									},
								}},
							{PC: 2,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r3"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r3"}},
										2,
									},
								}},
							{PC: 3,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r4"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r4"}},
										3,
									},
								}},
							{PC: 4,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r5"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r5"}},
										4,
									},
								}},
							{PC: 5,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r6"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r6"}},
										5,
									},
								}},
							{PC: 6, Type: executor.ObsTypeRollback, Value: 0},
							{PC: 7,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r8"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r8"}},
										7,
									},
								}},
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
				{ // r1 != 0のとき
					PC: 8,
					Registers: map[string]interface{}{
						"r2": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}},
								1,
							},
						},
						"r3": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r3"}},
								2,
							},
						},
						"r4": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r4"}},
								3,
							},
						},
						"r5": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r5"}},
								4,
							},
						},
						"r6": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r6"}},
								5,
							},
						},
						"r7": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r7"}},
								6,
							},
						},
						"r8": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r8"}},
								7,
							},
						},
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStart, Value: 0},
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{
								Op: "==",
								Operands: []interface{}{
									executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}},
									0,
								},
							}},
							{PC: 7,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r8"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r8"}},
										7,
									},
								}},
							{PC: 8, Type: executor.ObsTypeRollback, Value: 0},
							{PC: 1,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r2"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}},
										1,
									},
								}},
							{PC: 2,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r3"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r3"}},
										2,
									},
								}},
							{PC: 3,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r4"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r4"}},
										3,
									},
								}},
							{PC: 4,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r5"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r5"}},
										4,
									},
								}},
							{PC: 5,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r6"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r6"}},
										5,
									},
								}},
							{PC: 6,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r7"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r7"}},
										6,
									},
								}},
							{PC: 7,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r8"}},
								Value: executor.SymbolicExpr{
									Op: "+",
									Operands: []interface{}{
										executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r8"}},
										7,
									},
								}},
						},
						PathCond: executor.SymbolicExpr{
							Op: "!=",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}},
								0,
							},
						},
					},
				},
			},
			ExpectError: false,
		},

		{
			Name: "Nested if-else with symbolic conditions",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"r1", "3"}},      // 条件: r1 == 0 (最初の分岐)
				{Mnemonic: "add", Operands: []string{"r2", "r2", "1"}}, // 偽側: r2 = r2 + 1
				{Mnemonic: "jmp", Operands: []string{"7"}},             // 偽側の終了
				{Mnemonic: "beqz", Operands: []string{"r2", "6"}},      // 真側: 条件: r2 == 0
				{Mnemonic: "add", Operands: []string{"r3", "r3", "1"}}, // 真側の偽側: r3 = r3 + 1
				{Mnemonic: "jmp", Operands: []string{"7"}},             // 真側の偽側の終了
				{Mnemonic: "add", Operands: []string{"r4", "r4", "1"}}, // 真側の真側: r4 = r4 + 1
			},
			InitialConfig: &executor.Configuration{
				Registers: map[string]interface{}{
					"r3": 0,
					"r4": 0,
				},
				PC:     0,
				Memory: map[int]interface{}{},
				Trace:  executor.Trace{},
			},
			MaxSteps: 100,
			ExpectedConfigs: []executor.Configuration{
				{ // r1 == 0, r2 == 0 のとき
					PC: 7,
					Registers: map[string]interface{}{
						"r3": 0,
						"r4": 1,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStart, Value: 0},
							{PC: 0,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}}, 0}}},

							{PC: 1,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r2"}},
								Value:   executor.SymbolicExpr{Op: "+", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 1}}},
							{PC: 2,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
							{PC: 7,
								Type:  executor.ObsTypeRollback,
								Value: 0},
							{PC: 3,
								Type:  executor.ObsTypeStart,
								Value: 0}, //ここのvalueも通し番号にしたい？
							{PC: 3,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 0}}},
							{PC: 4,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r3"}},
								Value:   1},
							{PC: 5,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
							{PC: 7,
								Type:  executor.ObsTypeRollback,
								Value: 0},
							{PC: 6, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r4"}}, Value: 1},
						},
						PathCond: executor.SymbolicExpr{
							Op: "&&",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}}, 0}},
								executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 0}},
							},
						},
					},
				},
				{ // r1 == 0, r2 != 0 のとき
					PC: 7,
					Registers: map[string]interface{}{
						"r3": 1,
						"r4": 0,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStart, Value: 0},
							{PC: 0,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}}, 0}}},

							{PC: 1,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r2"}},
								Value:   executor.SymbolicExpr{Op: "+", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 1}}},
							{PC: 2,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
							{PC: 7,
								Type:  executor.ObsTypeRollback,
								Value: 0},
							{PC: 3,
								Type:  executor.ObsTypeStart,
								Value: 0}, //ここのvalueも通し番号にしたい？
							{PC: 3,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 0}}},
							{PC: 6, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r4"}}, Value: 1},
							{PC: 7,
								Type:  executor.ObsTypeRollback,
								Value: 0},
							{PC: 4,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r3"}},
								Value:   1},
							{PC: 5,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
						},
						PathCond: executor.SymbolicExpr{
							Op: "&&",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}}, 0}},
								executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 0}},
							},
						},
					},
				},
				{ // r1 != 0, r2 == 0 のとき
					PC: 7,
					Registers: map[string]interface{}{
						"r2": executor.SymbolicExpr{
							Op:       "+",
							Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 1},
						},
						"r3": 0,
						"r4": 0,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStart, Value: 0},
							{PC: 0,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}}, 0}}},
							{PC: 3,
								Type:  executor.ObsTypeStart,
								Value: 1}, //ここのvalueも通し番号にしたい？
							{PC: 3,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 0}}},
							{PC: 4,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r3"}},
								Value:   1},
							{PC: 5,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
							{PC: 7,
								Type:  executor.ObsTypeRollback,
								Value: 1},
							{PC: 6,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r4"}},
								Value:   1},
							{PC: 7,
								Type:  executor.ObsTypeRollback,
								Value: 0},
							{PC: 1,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r2"}},
								Value:   executor.SymbolicExpr{Op: "+", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 1}}},
							{PC: 2,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
						},
						PathCond: executor.SymbolicExpr{
							Op: "&&",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}}, 0}},
								executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 0}},
							},
						},
					},
				},
				{ // r1 != 0, r2 != 0 のとき
					PC: 7,
					Registers: map[string]interface{}{
						"r2": executor.SymbolicExpr{
							Op:       "+",
							Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 1},
						},
						"r3": 0,
						"r4": 0,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStart, Value: 0},
							{PC: 0,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}}, 0}}},
							{PC: 3,
								Type:  executor.ObsTypeStart,
								Value: 1}, //ここのvalueも通し番号にしたい？
							{PC: 3,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 0}}},
							{PC: 6,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r4"}},
								Value:   1},
							{PC: 7,
								Type:  executor.ObsTypeRollback,
								Value: 1},
							{PC: 4,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r3"}},
								Value:   1},
							{PC: 5,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
							{PC: 7,
								Type:  executor.ObsTypeRollback,
								Value: 0},
							{PC: 1,
								Type:    executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"r2"}},
								Value:   executor.SymbolicExpr{Op: "+", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 1}}},
							{PC: 2,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
						},
						PathCond: executor.SymbolicExpr{
							Op: "&&",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r1"}}, 0}},
								executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"r2"}}, 0}},
							},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "concreat single benq",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"r1", "10"}},
				{Mnemonic: "add", Operands: []string{"r2", "r2", "1"}},
			},
			InitialConfig: &executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"r1": 42,
				},
				Memory: map[int]interface{}{},
				Trace:  executor.Trace{},
			},
			MaxSteps: 10,
			ExpectedConfigs: []executor.Configuration{
				{
					PC: 2,
					Registers: map[string]interface{}{
						"r1": 42,
						"r2": executor.SymbolicExpr{
							Op: "+",
							Operands: []interface{}{executor.SymbolicExpr{
								Op:       "symbol",
								Operands: []interface{}{"r2"}}, 1}},
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypeStart, Value: 0},
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{42, 0}}},
							{PC: 10, Type: executor.ObsTypeRollback, Value: 0},
							{PC: 1,
								Type: executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{
									Op:       "var",
									Operands: []interface{}{"r2"}},
								Value: executor.SymbolicExpr{Op: "+", Operands: []interface{}{executor.SymbolicExpr{
									Op:       "symbol",
									Operands: []interface{}{"r2"}}, 1}}},
						},
						PathCond: executor.SymbolicExpr{
							Op:       "!=",
							Operands: []interface{}{42, 0},
						},
					},
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
	RunTestCase(t, testCases, specExecuteWrapper)
}

func specExecuteWrapper(program []assembler.OpCode, initialConfig *executor.Configuration, maxSteps int) ([]*executor.Configuration, error) {
	// ここで remainingWindow を固定値（例: 10）として渡す
	const remainingWindow = 5
	return executor.SpecExecute(program, initialConfig, maxSteps, remainingWindow)
}
