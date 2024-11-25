package executor_test

import (
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

func TestAlwaysMispredictStep(t *testing.T) {
	testCases := []struct {
		name           string
		initialConf    executor.Configuration
		instruction    assembler.OpCode
		expectedConfs  []executor.Configuration
		expectedIsSpec bool
		expectError    bool
	}{
		{
			name: "beqz concrete true",
			initialConf: executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"reg": 0,
				},
				Memory: map[int]interface{}{},
				Trace: executor.Trace{
					PathCond: executor.SymbolicExpr{},
				},
			},
			instruction: assembler.OpCode{
				Mnemonic: "beqz",
				Operands: []string{"reg", "4"},
			},
			expectedConfs: []executor.Configuration{
				{
					PC: 1,
					Registers: map[string]interface{}{
						"reg": 0,
					},
					Trace: executor.Trace{
						PathCond: executor.SymbolicExpr{
							Op: "==",
							Operands: []interface{}{
								"reg", 0,
							},
						},
						Observations: []executor.Observation{
							{
								PC:   0,
								Type: executor.ObsTypePC,
								Value: executor.SymbolicExpr{
									Op: "==",
									Operands: []interface{}{
										"reg", 0,
									},
								},
							},
						},
					},
				},
			},
			expectedIsSpec: true,
			expectError:    false,
		},
		{
			name: "beqz symbolic condition",
			initialConf: executor.Configuration{
				PC: 0,
				Registers: map[string]interface{}{
					"reg": executor.SymbolicExpr{
						Op:       ">",
						Operands: []interface{}{"x", 0},
					},
				},
				Memory: map[int]interface{}{},
				Trace: executor.Trace{
					PathCond: executor.SymbolicExpr{},
				},
			},
			instruction: assembler.OpCode{
				Mnemonic: "beqz",
				Operands: []string{"reg", "4"},
			},
			expectedConfs: []executor.Configuration{
				{
					PC: 1,
					Registers: map[string]interface{}{
						"reg": executor.SymbolicExpr{
							Op:       ">",
							Operands: []interface{}{"x", 0},
						},
					},
					Trace: executor.Trace{
						PathCond: executor.SymbolicExpr{
							Op: "==",
							Operands: []interface{}{
								"reg", 0,
							},
						},
						Observations: []executor.Observation{
							{
								PC:   0,
								Type: executor.ObsTypePC,
								Value: executor.SymbolicExpr{
									Op: "==",
									Operands: []interface{}{
										"reg", 0,
									},
								},
							},
						},
					},
				},
				{
					PC: 4,
					Registers: map[string]interface{}{
						"reg": executor.SymbolicExpr{
							Op:       ">",
							Operands: []interface{}{"x", 0},
						},
					},
					Trace: executor.Trace{
						PathCond: executor.SymbolicExpr{
							Op: "!=",
							Operands: []interface{}{
								"reg", 0,
							},
						},
						Observations: []executor.Observation{
							{
								PC:   0,
								Type: executor.ObsTypePC,
								Value: executor.SymbolicExpr{
									Op: "==",
									Operands: []interface{}{
										"reg", 0,
									},
								},
							},
						},
					},
				},
			},
			expectedIsSpec: true,
			expectError:    false,
		},
		{
			name: "jmp instruction",
			initialConf: executor.Configuration{
				PC:     0,
				Memory: map[int]interface{}{},
				Trace:  executor.Trace{},
			},
			instruction: assembler.OpCode{
				Mnemonic: "jmp",
				Operands: []string{"4"},
			},
			expectedConfs: []executor.Configuration{
				{
					PC: 4,
					Trace: executor.Trace{
						Observations: []executor.Observation{
							{
								PC:    0,
								Type:  executor.ObsTypePC,
								Value: executor.SymbolicExpr{Op: "jmp", Operands: []interface{}{4}},
							},
						},
					},
				},
			},
			expectedIsSpec: false,
			expectError:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			confs, isSpec, err := executor.AlwaysMispredictStep(testCase.instruction, &testCase.initialConf)

			if (err != nil) != testCase.expectError {
				t.Errorf("Test case '%s' failed: expected an error but got none", testCase.name)

				return
			}
			if len(confs) != len(testCase.expectedConfs) {
				t.Errorf("FAIL: %s - expected %d configurations, got %d\n", testCase.name, len(testCase.expectedConfs), len(confs))

				return
			}
			for i, conf := range confs {
				if !executor.CompareConfiguration(testCase.expectedConfs[i], *conf) {
					differences := executor.FormatConfigDifferences(testCase.expectedConfs[i], *conf)
					t.Errorf("Test case '%s' failed: Trace %d did not match expected confs.\n%s",
						testCase.name, i+1, differences)
				}
			}
			if isSpec != testCase.expectedIsSpec {
				t.Errorf("FAIL: %s - expected speculative: %v, got: %v\n", testCase.name, testCase.expectedIsSpec, isSpec)
			}
		})
	}
}
