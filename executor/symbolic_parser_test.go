package executor

import (
	"reflect"
	"testing"
)

func TestParseSymbolicExpr(t *testing.T) {
	testCases := []struct {
		input                string
		expectedSymbolicExpr *SymbolicExpr
		expectedError        bool
	}{
		// 正常系: 単純な加算
		{
			input: "10 + 20",
			expectedSymbolicExpr: &SymbolicExpr{
				Op: "+",
				Operands: []interface{}{
					SymbolicExpr{
						Op: "value",
						Operands: []interface{}{
							10,
						},
					},
					SymbolicExpr{
						Op: "value",
						Operands: []interface{}{
							20,
						},
					},
				},
			},
			expectedError: false,
		},
		// 正常系: ネストした式
		{
			input: "10 + x * 2",
			expectedSymbolicExpr: &SymbolicExpr{
				Op: "+",
				Operands: []interface{}{
					SymbolicExpr{
						Op: "value",
						Operands: []interface{}{
							10,
						},
					},
					SymbolicExpr{
						Op: "*",
						Operands: []interface{}{
							SymbolicExpr{
								Op: "value",
								Operands: []interface{}{
									"x",
								},
							},
							SymbolicExpr{
								Op: "value",
								Operands: []interface{}{
									2,
								},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		// 正常系: 比較演算子
		{
			input: "x > 10",
			expectedSymbolicExpr: &SymbolicExpr{
				Op: ">",
				Operands: []interface{}{
					SymbolicExpr{
						Op: "value",
						Operands: []interface{}{
							"x",
						},
					},
					SymbolicExpr{
						Op: "value",
						Operands: []interface{}{
							10,
						},
					},
				},
			},
			expectedError: false,
		},
		// 正常系: 複数の演算子
		{
			input: "(x + 10) * 20",
			expectedSymbolicExpr: &SymbolicExpr{
				Op: "*",
				Operands: []interface{}{
					SymbolicExpr{
						Op: "+",
						Operands: []interface{}{
							SymbolicExpr{
								Op: "value",
								Operands: []interface{}{
									"x",
								},
							},
							SymbolicExpr{
								Op: "value",
								Operands: []interface{}{
									10,
								},
							},
						},
					},
					SymbolicExpr{
						Op: "value",
						Operands: []interface{}{
							20,
						},
					},
				},
			},
			expectedError: false,
		},
		// 異常系: 不正なトークン
		{
			input:                "x + @",
			expectedSymbolicExpr: nil,
			expectedError:        true,
		},
		// 異常系: 括弧が一致しない
		{
			input:                "(x + 10",
			expectedSymbolicExpr: nil,
			expectedError:        true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			result, err := ParseSymbolicExpr(testCase.input)

			if testCase.expectedError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !reflect.DeepEqual(result, testCase.expectedSymbolicExpr) {
					t.Errorf("expected %+v, got %+v", testCase.expectedSymbolicExpr, result)
				}
			}
		})
	}
}
