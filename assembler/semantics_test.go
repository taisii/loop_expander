package assembler_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/taisii/go-project/assembler"
)

func TestStep(t *testing.T) {
	tests := []struct {
		name           string
		instruction    assembler.OpCode
		initialConfig  *assembler.Configuration
		expectedConfig *assembler.Configuration
		expectedError  error
	}{
		{
			name: "MOV instruction",
			instruction: assembler.OpCode{
				Mnemonic: "mov",
				Operands: []string{"r1", "10"},
			},
			initialConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 0},
				PC:        0,
			},
			expectedConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 10},
				PC:        1,
			},
			expectedError: nil,
		},
		{
			name: "ADD instruction",
			instruction: assembler.OpCode{
				Mnemonic: "add",
				Operands: []string{"r2", "r1", "5"},
			},
			initialConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 10, "r2": 0},
				PC:        0,
			},
			expectedConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 10, "r2": 15},
				PC:        1,
			},
			expectedError: nil,
		},
		{
			name: "BEQZ instruction with zero",
			instruction: assembler.OpCode{
				Mnemonic: "beqz",
				Operands: []string{"r1", "3"},
			},
			initialConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 0},
				PC:        0,
			},
			expectedConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 0},
				PC:        3,
			},
			expectedError: nil,
		},
		{
			name: "BEQZ instruction with non-zero",
			instruction: assembler.OpCode{
				Mnemonic: "beqz",
				Operands: []string{"r1", "3"},
			},
			initialConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 1},
				PC:        0,
			},
			expectedConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 1},
				PC:        1,
			},
			expectedError: nil,
		},
		{
			name: "Unsupported instruction",
			instruction: assembler.OpCode{
				Mnemonic: "unknown",
				Operands: []string{},
			},
			initialConfig: &assembler.Configuration{
				Registers: map[string]int{},
				PC:        0,
			},
			expectedConfig: nil,
			expectedError:  errors.New("unsupported instruction: unknown"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config, err := assembler.Step(test.instruction, test.initialConfig)

			// Config の比較
			if test.expectedConfig != nil && config != nil {
				if config.PC != test.expectedConfig.PC || !equalRegisters(config.Registers, test.expectedConfig.Registers) {
					t.Errorf("Expected configuration: %+v, got: %+v", test.expectedConfig, config)
				}
			} else if test.expectedConfig == nil && config != nil {
				t.Errorf("Expected nil configuration, got: %+v", config)
			}

			// Error の比較
			if test.expectedError == nil && err != nil {
				t.Errorf("Unexpected error: %v", err)
			} else if test.expectedError != nil && err == nil {
				t.Errorf("Expected error: %v, got nil", test.expectedError)
			} else if test.expectedError != nil && err != nil {
				// エラーメッセージが等しいかどうかを比較
				if err.Error() != test.expectedError.Error() {
					t.Errorf("Expected error: %v, got: %v", test.expectedError, err)
				}
			}
		})
	}
}

// equalRegisters compares two maps of registers for equality
func equalRegisters(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func TestRunProgram(t *testing.T) {
	// 正常に終了するプログラムセット
	successProgram := []assembler.OpCode{
		{Mnemonic: "mov", Operands: []string{"r1", "10"}},      // r1 = 10
		{Mnemonic: "add", Operands: []string{"r2", "r1", "5"}}, // r2 = r1 + 5
		{Mnemonic: "beqz", Operands: []string{"r2", "5"}},      // if r2 == 0, jump to PC = 5
		{Mnemonic: "jmp", Operands: []string{"4"}},             // jump to PC = 4 (end point)
	}

	// タイムアウトするプログラムセット（無限ループ）
	timeoutProgram := []assembler.OpCode{
		{Mnemonic: "mov", Operands: []string{"r1", "0"}},       // r1 = 0
		{Mnemonic: "add", Operands: []string{"r1", "r1", "1"}}, // r1 = r1 + 1
		{Mnemonic: "jmp", Operands: []string{"1"}},             // jump to PC = 1 (loop)
	}

	tests := []struct {
		name           string
		program        []assembler.OpCode
		maxSteps       int
		initialConfig  *assembler.Configuration
		expectedConfig *assembler.Configuration
		expectedError  error
	}{
		{
			name:     "Normal execution",
			program:  successProgram,
			maxSteps: 10,
			initialConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 0, "r2": 0},
				PC:        0,
			},
			expectedConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 10, "r2": 15},
				PC:        4, // Program ends here
			},
			expectedError: nil,
		},
		{
			name:     "Timeout due to infinite loop",
			program:  timeoutProgram,
			maxSteps: 10,
			initialConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 0},
				PC:        0,
			},
			expectedConfig: &assembler.Configuration{
				Registers: map[string]int{"r1": 10}, // After 10 steps
				PC:        1,                        // Stuck in the loop
			},
			expectedError: errors.New("timeout reached"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			finalConfig, err := assembler.Run(test.program, test.initialConfig, test.maxSteps)

			// エラー比較を最優先
			if (err != nil || test.expectedError != nil) && (err == nil || test.expectedError == nil || err.Error() != test.expectedError.Error()) {
				t.Errorf("Expected error: %v, got: %v", test.expectedError, err)
				return // エラーが一致しない場合はここで終了
			}

			// タイムアウトではない場合のみ、コンフィグを比較
			if err == nil {
				if !reflect.DeepEqual(finalConfig.Registers, test.expectedConfig.Registers) || finalConfig.PC != test.expectedConfig.PC {
					t.Errorf("Expected configuration: %+v, got: %+v", test.expectedConfig, finalConfig)
				}
			}
		})
	}
}
