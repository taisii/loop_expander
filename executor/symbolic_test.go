package executor

import (
	"testing"
)

func TestNewConfiguration(t *testing.T) {
	memory := map[int]interface{}{0: 42, 1: 84}
	registers := map[string]interface{}{"r1": 10, "r2": 20}
	conf := NewConfiguration(memory, registers)

	if conf.PC != 0 {
		t.Errorf("Expected PC to be 0, got %d", conf.PC)
	}

	if conf.Memory[0] != 42 || conf.Memory[1] != 84 {
		t.Errorf("Unexpected memory values: %v", conf.Memory)
	}

	if conf.Registers["r1"] != 10 || conf.Registers["r2"] != 20 {
		t.Errorf("Unexpected register values: %v", conf.Registers)
	}
}

func TestEvalExprConcreteValues(t *testing.T) {
	conf := NewConfiguration(nil, map[string]interface{}{"r1": 10})

	value, err := evalExpr(5, conf)
	if err != nil || value != 5 {
		t.Errorf("Expected 5, got %v (err: %v)", value, err)
	}

	value, err = evalExpr("r1", conf)
	if err != nil || value != 10 {
		t.Errorf("Expected 10, got %v (err: %v)", value, err)
	}

	value, err = evalExpr("unknown", conf)
	expectedValue := SymbolicExpr{
		Op:       "symbol",
		Operands: []interface{}{"unknown"},
	}
	if err != nil || !CompareSymbolicExpr(value, expectedValue) {
		t.Errorf("Expected symbolic value for unknown register, got %v", value)
	}
}

// Test case structure
type EvalTestCase struct {
	name           string         // Test case name
	initialConf    *Configuration // Initial configuration
	symbolicExpr   interface{}    // Expression to evaluate
	expectedResult interface{}    // Expected result
	expectError    bool           // Whether an error is expected
}

func TestEvalExprWithStruct(t *testing.T) {
	// Define test cases
	testCases := []EvalTestCase{
		{
			name: "Concrete value",
			initialConf: NewConfiguration(
				nil,
				map[string]interface{}{},
			),
			symbolicExpr:   42,
			expectedResult: 42,
			expectError:    false,
		},
		{
			name: "Register value",
			initialConf: NewConfiguration(
				nil,
				map[string]interface{}{
					"r1": 10,
				},
			),
			symbolicExpr:   "r1",
			expectedResult: 10,
			expectError:    false,
		},
		{
			name: "Evaluatable symbolic expression",
			initialConf: NewConfiguration(
				nil,
				map[string]interface{}{
					"r1": 10,
					"r2": 20,
				},
			),
			symbolicExpr: SymbolicExpr{
				Op:       "+",
				Operands: []interface{}{"r1", "r2"},
			},
			expectedResult: 30,
			expectError:    false,
		},
		{
			name: "Partially evaluatable symbolic expression",
			initialConf: NewConfiguration(
				nil,
				map[string]interface{}{
					"r1": 10,
				},
			),
			symbolicExpr: SymbolicExpr{
				Op:       "+",
				Operands: []interface{}{"r1", "unknown"},
			},
			expectedResult: SymbolicExpr{
				Op: "+",
				Operands: []interface{}{10, SymbolicExpr{
					Op:       "symbol",
					Operands: []interface{}{"unknown"},
				}},
			},
			expectError: false,
		},
		{
			name: "Unsupported expression type",
			initialConf: NewConfiguration(
				nil,
				map[string]interface{}{},
			),
			symbolicExpr:   map[string]interface{}{"invalid": "expression"},
			expectedResult: nil,
			expectError:    true,
		},
		{
			name: "nested greater than",
			initialConf: NewConfiguration(
				nil,
				map[string]interface{}{
					"x": SymbolicExpr{
						Op:       ">",
						Operands: []interface{}{10, "y"},
					},
					"y": 5,
				},
			),
			symbolicExpr: SymbolicExpr{
				Op:       ">",
				Operands: []interface{}{10, "x"}, // x > y
			},
			expectedResult: 1, // True (concrete evaluation)
			expectError:    false,
		},
		{
			name: "Symbolic condition: less than with unknown",
			initialConf: NewConfiguration(
				nil,
				map[string]interface{}{
					"x": 15,
				},
			),
			symbolicExpr: SymbolicExpr{
				Op:       "<",
				Operands: []interface{}{"x", "unknown"}, // x < unknown
			},
			expectedResult: SymbolicExpr{
				Op: "<",
				Operands: []interface{}{15, SymbolicExpr{
					Op:       "symbol",
					Operands: []interface{}{"unknown"},
				}},
			}, // Symbolic result
			expectError: false,
		},
		{
			name: "Symbolic condition: equality with concrete",
			initialConf: NewConfiguration(
				nil,
				map[string]interface{}{
					"a": 20,
					"b": 20,
				},
			),
			symbolicExpr: SymbolicExpr{
				Op:       "==",
				Operands: []interface{}{"a", "b"}, // a == b
			},
			expectedResult: 1, // True (concrete evaluation)
			expectError:    false,
		},
		{
			name: "Symbolic condition: inequality with symbolic operand",
			initialConf: NewConfiguration(
				nil,
				map[string]interface{}{
					"a": 10,
				},
			),
			symbolicExpr: SymbolicExpr{
				Op:       "!=",
				Operands: []interface{}{"a", "unknown"}, // a != unknown
			},
			expectedResult: SymbolicExpr{
				Op: "!=",
				Operands: []interface{}{10, SymbolicExpr{
					Op:       "symbol",
					Operands: []interface{}{"unknown"},
				}},
			}, // Symbolic result
			expectError: false,
		},
		{
			name: "Nested symbolic concreat value",
			initialConf: &Configuration{
				Registers: map[string]interface{}{
					"x": 2,
					"y": 0,
				},
			},
			symbolicExpr: SymbolicExpr{
				Op: "+",
				Operands: []interface{}{10, SymbolicExpr{
					Op:       "*",
					Operands: []interface{}{2, "x"},
				}},
			},
			expectedResult: 14,
			expectError:    false,
		},
		{
			name: "Nested parsed concreat value",
			initialConf: &Configuration{
				Registers: map[string]interface{}{
					"x": 2,
					"y": 0,
				},
			},
			symbolicExpr: SymbolicExpr{
				Op: "+",
				Operands: []interface{}{10, SymbolicExpr{
					Op: "*",
					Operands: []interface{}{SymbolicExpr{
						Op:       "value",
						Operands: []interface{}{2},
					}, SymbolicExpr{
						Op:       "value",
						Operands: []interface{}{"x"},
					}},
				}},
			},
			expectedResult: 14,
			expectError:    false,
		},
	}

	// Execute each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := evalExpr(tc.symbolicExpr, tc.initialConf)
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !CompareSymbolicExpr(result, tc.expectedResult) {
				t.Errorf("Expected result: %v, got: %v", tc.expectedResult, result)
			}
		})
	}
}

func TestAllConcrete(t *testing.T) {
	operands := []interface{}{1, 2, 3}
	if !allConcrete(operands) {
		t.Errorf("Expected all concrete, got false")
	}

	operands = []interface{}{1, "r1", 3}
	if allConcrete(operands) {
		t.Errorf("Expected not all concrete, got true")
	}
}

func TestComputeConcrete(t *testing.T) {
	result, err := computeConcrete("+", []interface{}{3, 7})
	if err != nil || result != 10 {
		t.Errorf("Expected 10, got %v (err: %v)", result, err)
	}

	result, err = computeConcrete("-", []interface{}{10, 5})
	if err != nil || result != 5 {
		t.Errorf("Expected 5, got %v (err: %v)", result, err)
	}

	result, err = computeConcrete("*", []interface{}{2, 4})
	if err != nil || result != 8 {
		t.Errorf("Expected 8, got %v (err: %v)", result, err)
	}

	result, err = computeConcrete("/", []interface{}{8, 2})
	if err != nil || result != 4 {
		t.Errorf("Expected 4, got %v (err: %v)", result, err)
	}

	result, err = computeConcrete("/", []interface{}{8, 0})
	if err == nil {
		t.Errorf("Expected division by zero error, got result: %v", result)
	}

	_, err = computeConcrete("invalid_op", []interface{}{1, 2})
	if err == nil {
		t.Errorf("Expected error for invalid operator")
	}
}
