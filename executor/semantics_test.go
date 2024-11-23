package executor_test

import (
	"encoding/json"
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

func TestStep(t *testing.T) {
	tests := []struct {
		name           string
		initialConf    *executor.Configuration
		program        []assembler.OpCode
		expectedPC     int
		expectedStates []map[string]interface{}
		expectError    bool
	}{
		{
			name: "Symbolic addition",
			initialConf: executor.NewConfiguration(
				map[int]interface{}{},
				map[string]interface{}{
					"x": executor.SymbolicExpr{
						Op:       "+",
						Operands: []interface{}{5, "y"},
					},
					"y": 3,
				}),
			program: []assembler.OpCode{
				{Mnemonic: "add", Operands: []string{"z", "x", "y"}}, // z = (5 + y) + y
			},
			expectedPC: 1,
			expectedStates: []map[string]interface{}{
				{"x": executor.SymbolicExpr{
					Op:       "+",
					Operands: []interface{}{5, "y"},
				},
					"y": 3, "z": 11},
			},
			expectError: false,
		},
		{
			name: "Conditional branch with symbolic register",
			initialConf: executor.NewConfiguration(
				map[int]interface{}{},
				map[string]interface{}{
					"x": executor.SymbolicExpr{
						Op:       ">",
						Operands: []interface{}{10, "y"},
					},
					"y": 5,
				}),
			program: []assembler.OpCode{
				{Mnemonic: "beqz", Operands: []string{"x", "100"}}, // x > y evaluates symbolically
			},
			expectedPC: 1, // PC should increment if branch condition is false
			expectedStates: []map[string]interface{}{
				{"x": executor.SymbolicExpr{
					Op:       ">",
					Operands: []interface{}{10, "y"},
				}},
			},
			expectError: false,
		},
		{
			name: "Symbolic jump",
			initialConf: executor.NewConfiguration(
				map[int]interface{}{},
				map[string]interface{}{
					"x": executor.SymbolicExpr{
						Op:       "*",
						Operands: []interface{}{2, "y"},
					},
					"y": 4,
				}),
			program: []assembler.OpCode{
				{Mnemonic: "jmp", Operands: []string{"x"}}, // Jump to PC = 2 * y
			},
			expectedPC: 8, // PC = 2 * 4
			expectedStates: []map[string]interface{}{
				{"x": executor.SymbolicExpr{
					Op:       "*",
					Operands: []interface{}{2, "y"},
				}},
			},
			expectError: false,
		},
		{
			name: "Unsupported instruction error",
			initialConf: executor.NewConfiguration(
				map[int]interface{}{},
				map[string]interface{}{"x": 1}),
			program: []assembler.OpCode{
				{Mnemonic: "unknown", Operands: []string{"x", "y"}},
			},
			expectedPC:     0, // PC remains unchanged
			expectedStates: []map[string]interface{}{},
			expectError:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conf := test.initialConf

			_, err := executor.ExecuteProgram(test.program, conf, 10)

			if test.expectError {
				if err == nil {
					t.Errorf("expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect an error but got: %v", err)
				}

				// Check PC value
				if conf.PC != test.expectedPC {
					t.Errorf("expected PC: %d, got: %d", test.expectedPC, conf.PC)
				}

				// Check register states
				for _, state := range test.expectedStates {
					for key, value := range state {
						if !compareValues(conf.Registers[key], value) {
							t.Errorf("for register '%s', expected value: %v, got: %v", key, value, conf.Registers[key])
						}
					}
				}
			}
		})
	}
}

func compareValues(a, b interface{}) bool {
	jsonA, errA := json.Marshal(a)
	jsonB, errB := json.Marshal(b)
	if errA != nil || errB != nil {
		return false
	}
	return string(jsonA) == string(jsonB)
}
