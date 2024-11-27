package executor_test

import (
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

func TestExecute(t *testing.T) {
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
			},
			MaxSteps: 10,
			ExpectedConfigs: []executor.Configuration{
				{
					PC:        4,
					StepCount: 2,
					Registers: map[string]interface{}{
						"y": 2,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{
								Op: "symbol", Operands: []interface{}{"x"}}, 0}}},
							{PC: 3,
								Type: executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{
									Op:       "var",
									Operands: []interface{}{"y"}},
								Value: 2},
						},
						PathCond: executor.SymbolicExpr{
							Op: "==",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"x"}},
								0,
							},
						},
					},
				},
				{
					PC:        5,
					StepCount: 3,
					Registers: map[string]interface{}{
						"y": 1,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{
								Op: "symbol", Operands: []interface{}{"x"}}, 0}}},
							{PC: 1,
								Type: executor.ObsTypeStore,
								Address: &executor.SymbolicExpr{
									Op:       "var",
									Operands: []interface{}{"y"}},
								Value: 1},
							{PC: 2, Type: executor.ObsTypePC, Address: nil, Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{5}}},
						},
						PathCond: executor.SymbolicExpr{
							Op: "!=",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"x"}},
								0,
							},
						},
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
				Trace:     executor.Trace{},
			},
			MaxSteps:        5,
			ExpectedConfigs: nil,
			ExpectError:     false,
		},
		{
			Name: "Nested If Statements",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"x", "3"}},
				{Mnemonic: "mov", Operands: []string{"y", "1"}},
				{Mnemonic: "jmp", Operands: []string{"7"}},
				{Mnemonic: "beqz", Operands: []string{"y", "6"}},
				{Mnemonic: "mov", Operands: []string{"z", "2"}},
				{Mnemonic: "jmp", Operands: []string{"7"}},
				{Mnemonic: "mov", Operands: []string{"z", "3"}},
			},
			InitialConfig: &executor.Configuration{
				Registers: make(map[string]interface{}),
				PC:        0,
				Trace:     executor.Trace{},
			},
			MaxSteps: 10,
			ExpectedConfigs: []executor.Configuration{
				{
					PC:        7,
					StepCount: 3,
					Registers: map[string]interface{}{
						"z": 3,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"x"}}, 0}}},
							{PC: 3, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"y"}}, 0}}},
							{PC: 6, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"z"}}, Value: 3},
						},
						PathCond: executor.SymbolicExpr{
							Op: "&&",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"x"}}, 0}},
								executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"y"}}, 0}},
							},
						},
					},
				},
				{
					PC:        7,
					StepCount: 3,
					Registers: map[string]interface{}{
						"y": 1,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"x"}}, 0}}},
							{PC: 1, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"y"}}, Value: 1},
							{PC: 2, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
						},
						PathCond: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"x"}}, 0}},
					},
				},
				{
					PC:        7,
					StepCount: 4,
					Registers: map[string]interface{}{
						"z": 2,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"x"}}, 0}}},
							{PC: 3, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"y"}}, 0}}},
							{PC: 4, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"z"}}, Value: 2},
							{PC: 5, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{7}}},
						},
						PathCond: executor.SymbolicExpr{
							Op: "&&",
							Operands: []interface{}{
								executor.SymbolicExpr{Op: "==", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"x"}}, 0}},
								executor.SymbolicExpr{Op: "!=", Operands: []interface{}{executor.SymbolicExpr{Op: "symbol", Operands: []interface{}{"y"}}, 0}},
							},
						},
					},
				},
			},
			ExpectError: false,
		},
		{
			Name: "Concreat One Branch",
			Program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"x", "3"}},
				{Mnemonic: "mov", Operands: []string{"y", "1"}},
				{Mnemonic: "jmp", Operands: []string{"0"}},
				{Mnemonic: "mov", Operands: []string{"z", "2"}},
			},
			InitialConfig: &executor.Configuration{
				Registers: map[string]interface{}{
					"x": 0,
				},
				PC:    0,
				Trace: executor.Trace{},
			},
			MaxSteps: 10,
			ExpectedConfigs: []executor.Configuration{
				{
					PC:        4,
					StepCount: 2,
					Registers: map[string]interface{}{
						"x": 0,
						"z": 2,
					},
					Memory: map[int]interface{}{},
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{PC: 0, Type: executor.ObsTypePC, Value: executor.SymbolicExpr{Op: "==", Operands: []interface{}{0, 0}}},
							{PC: 3, Type: executor.ObsTypeStore, Address: &executor.SymbolicExpr{Op: "var", Operands: []interface{}{"z"}}, Value: 2},
						},
						PathCond: executor.SymbolicExpr{
							Op:       "==",
							Operands: []interface{}{0, 0},
						},
					},
				},
			},
			ExpectError: false,
		},
	}

	// テストケースの実行
	RunTestCase(t, testCases, executor.ExecuteProgram)
}
