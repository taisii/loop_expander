package loop_expander_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/loop_expander"
)

type testCase struct {
	name           string
	inputAsm       *assembler.Assembler
	maxUnrollCount int
	expectedAsm    *assembler.Assembler
	expectedError  error
}

func TestLoop_expander(t *testing.T) {
	testCases := []testCase{
		{
			name: "No Loop", // ループがない場合
			inputAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
				},
				Labels: map[string]int{},
			},
			maxUnrollCount: 2,
			expectedAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
				},
				Labels: map[string]int{},
			},
		},
		{
			name: "Simple Loop", // 簡単なループ
			inputAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart"}}},
				},
				Labels: map[string]int{"LoopStart": 0},
			},
			maxUnrollCount: 3,
			expectedAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_0"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_1"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_2"}}},
					{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
				},
				Labels: map[string]int{
					"LoopStart":   0,
					"LoopStart_0": 3,
					"LoopStart_1": 6,
					"LoopStart_2": 9,
					"programEnd":  9},
			},
		},
		{
			name: "basic loop",
			inputAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart"}}},
				},
				Labels: map[string]int{
					"LoopStart": 1,
				},
			},
			maxUnrollCount: 3,
			expectedAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_0"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_1"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_2"}}},
					{Addr: 9, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
				},
				Labels: map[string]int{
					"LoopStart":   1,
					"LoopStart_0": 4,
					"LoopStart_1": 7,
					"LoopStart_2": 10,
					"programEnd":  10,
				},
			},
			expectedError: nil,
		},
		{
			name: "loop start from 0",
			inputAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart"}}},
				},
				Labels: map[string]int{
					"LoopStart": 0,
				},
			},
			maxUnrollCount: 3,
			expectedAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_0"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_0"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_1"}}},
					{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_1"}}},
					{Addr: 9, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 10, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 11, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_2"}}},
					{Addr: 12, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 13, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"LoopStart_2"}}},
					{Addr: 14, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
				},
				Labels: map[string]int{
					"LoopStart":   0,
					"LoopStart_0": 5,
					"LoopStart_1": 10,
					"LoopStart_2": 15,
					"programEnd":  15,
				},
			},
			expectedError: nil,
		},
		{
			name: "Loop with beqz",
			inputAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "LoopStart"}}},
				},
				Labels: map[string]int{
					"LoopStart": 1,
				},
			},
			maxUnrollCount: 3,
			expectedAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "LoopStart_0"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "LoopStart_1"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "LoopStart_2"}}},
					{Addr: 9, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
				},
				Labels: map[string]int{
					"LoopStart":   1,
					"LoopStart_0": 4,
					"LoopStart_1": 7,
					"LoopStart_2": 10,
					"programEnd":  10,
				},
			},
			expectedError: nil,
		},
		{
			name: "no loop with spbarr",
			inputAsm: &assembler.Assembler{
				Labels: map[string]int{
					"End": 5,
				},
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "<-", Operands: []string{"x", "v<y"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "End"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "spbarr", Operands: []string{""}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"v", "v"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"v", "v"}}},
				},
			},
			maxUnrollCount: 3,
			expectedAsm: &assembler.Assembler{
				Labels: map[string]int{
					"End": 5,
				},
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "<-", Operands: []string{"x", "v<y"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "End"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "spbarr", Operands: []string{""}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"v", "v"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"v", "v"}}},
				},
			},
			expectedError: nil,
		},
		{
			name: "Complex Loop with Multiple Internal Labels",
			inputAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "<-", Operands: []string{"w", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "<-", Operands: []string{"x", "in>=bound"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "L3"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"L10"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"secret", "in"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"z", "secret"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "Loop"}}},
				},
				Labels: map[string]int{
					"Loop": 1,
					"L3":   4,
					"L10":  6,
				},
			},
			maxUnrollCount: 2,
			expectedAsm: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "<-", Operands: []string{"w", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "<-", Operands: []string{"x", "in>=bound"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "L3"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"L10"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"secret", "in"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"z", "secret"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "Loop_0"}}},
					{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "<-", Operands: []string{"x", "in>=bound"}}},
					{Addr: 9, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "L3_0"}}},
					{Addr: 10, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"L10_0"}}},
					{Addr: 11, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"secret", "in"}}},
					{Addr: 12, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"z", "secret"}}},
					{Addr: 13, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "Loop_1"}}},
					{Addr: 14, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
				},
				Labels: map[string]int{
					"programEnd": 15,
					"L3":       4,
					"L3_0":       11,
					"L10":      6,
					"L10_0":      13,
					"Loop":       1,
					"Loop_0":     8,
					"Loop_1":     15,
					"L3_1":       18,
					"L10_1":      20,
				},
			},
			expectedError: nil,
		},
		// ネストされたループのテストケース、今回の論文では対応しない
		// {
		// 	name: "nested loop (not supported)",
		// 	inputAsm: &assembler.Assembler{
		// 		Program: []assembler.Instruction{
		// 			{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 			{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
		// 			{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"InnerLoop"}}},
		// 			{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
		// 			{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "InnerLoop"}}},
		// 			{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 		},
		// 		Labels: map[string]int{
		// 			"OuterLoop": 1,
		// 			"InnerLoop": 3,
		// 		},
		// 	},
		// 	maxUnrollCount: 2,
		// 	expectedAsm: &assembler.Assembler{
		// 		Program: []assembler.Instruction{
		// 			{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 			{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
		// 			{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"InnerLoop"}}},
		// 			{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
		// 			{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "InnerLoop_0"}}},
		// 			{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 			{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
		// 			{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "InnerLoop_1"}}},
		// 			{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 		},
		// 		Labels: map[string]int{
		// 			"OuterLoop":   1,
		// 			"InnerLoop":   3,
		// 			"InnerLoop_0": 6,
		// 			"InnerLoop_1": 9,
		// 		},
		// 	},
		// 	expectedError: nil,
		// },
		// ネストされたループのテストケース、今回の論文では対応しない
		// 		{
		// 	name: "loop_expander",
		// 	inputAsm: &assembler.Assembler{
		// 		Program: []assembler.Instruction{
		// 			{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 			{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
		// 			{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"InnerLoop"}}},
		// 			{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
		// 			{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "InnerLoop_0"}}},
		// 			{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 			{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
		// 			{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "InnerLoop_1"}}},
		// 			{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 		},
		// 		Labels: map[string]int{
		// 			"OuterLoop":   1,
		// 			"InnerLoop":   3,
		// 			"InnerLoop_0": 6,
		// 			"InnerLoop_1": 9,
		// 		},
		// 	},
		// 	maxUnrollCount: 2,
		// 	expectedAsm: &assembler.Assembler{
		// 		Program: []assembler.Instruction{
		// 			{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 			{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
		// 			{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"InnerLoop"}}},
		// 			{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
		// 			{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "InnerLoop_0"}}},
		// 			{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 			{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
		// 			{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "InnerLoop_1"}}},
		// 			{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
		// 		},
		// 		Labels: map[string]int{
		// 			"OuterLoop":   1,
		// 			"InnerLoop":   3,
		// 			"InnerLoop_0": 6,
		// 			"InnerLoop_1": 9,
		// 		},
		// 	},
		// 	expectedError: nil,
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resultAsm, err := loop_expander.Loop_expander(tc.inputAsm, tc.maxUnrollCount)

			if tc.expectedError != nil {
				if !errors.Is(err, tc.expectedError) {
					t.Errorf("expected error %v, but got %v", tc.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !assembler.CompareAssembler(resultAsm, tc.expectedAsm) {
				expectedCFG, errExpectedCFG := loop_expander.BuildControlFlowGraph(tc.expectedAsm)
				var expectedDot string
				if errExpectedCFG != nil {
					expectedDot = fmt.Sprintf("Error building expected CFG: %v", errExpectedCFG)
				} else if expectedCFG == nil {
					expectedDot = "Expected CFG is nil (no loop or empty program)"
				} else {
					expectedDot = loop_expander.ToDOT(expectedCFG)
				}

				actualCFG, errActualCFG := loop_expander.BuildControlFlowGraph(resultAsm)
				var actualDot string
				if errActualCFG != nil {
					actualDot = fmt.Sprintf("Error building actual CFG: %v", errActualCFG)
				} else if actualCFG == nil {
					actualDot = "Actual CFG is nil (no loop or empty program)"
				} else {
					actualDot = loop_expander.ToDOT(actualCFG)
				}

				t.Errorf("%s differs from expected.\n\nExpected Assembly:\n%s\nActual Assembly:\n%s\nAssembly Diff:\n%s\n\nExpected CFG (DOT):\n%s\nActual CFG (DOT):\n%s",
					tc.name,
					assembler.FormatAsm(tc.expectedAsm),
					assembler.FormatAsm(resultAsm),
					assembler.DiffAssembler(resultAsm, tc.expectedAsm),
					expectedDot,
					actualDot,
				)
			}
		})
	}
}
