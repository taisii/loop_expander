package loop_expander_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/loop_expander"
)

func TestDetectLoops(t *testing.T) {
	testCases := []struct {
		name     string
		assembly *assembler.Assembler
		expected [][]int
	}{
		{
			name: "No loops",
			assembly: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "L1"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"z", "2"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"w", "3"}}},
				},
				Labels: map[string]int{
					"L1": 4,
				},
			},
			expected: [][]int{}, // ループなし
		},
		{
			name: "Single loop",
			assembly: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "L1"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"z", "2"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"L2"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"w", "3"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"v", "4"}}},
				},
				Labels: map[string]int{
					"L1": 5,
					"L2": 2,
				},
			},
			expected: [][]int{
				{1, 2}, // 1 -> 2 -> 1 のループ
			},
		},
		{
			name: "Multiple loops",
			assembly: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "L1"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"z", "2"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"L2"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"w", "3"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "L3"}}},
					{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"v", "4"}}},
					{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"L4"}}},
				},
				Labels: map[string]int{
					"L1": 5,
					"L2": 2,
					"L3": 8,
					"L4": 5,
				},
			},
			expected: [][]int{
				{1, 2},    // 1 -> 2 -> 1 のループ
				{3, 4, 5}, // 3 -> 4 -> 5 -> 3 のループ
			},
		},
		{
			name: "Nested loops",
			assembly: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"L5"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"L5"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"z", "2"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"L1"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"w", "3"}}},
					{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"L1"}}},
				},
				Labels: map[string]int{
					"L1": 1,
					"L2": 2,
					"L3": 3,
					"L4": 4,
					"L5": 5,
					"L6": 6,
				},
			},
			expected: [][]int{ // ループが少なくとも一つ検出できていればよい
				{1, 2, 3, 4, 5},
				{1, 2, 3, 4, 5, 6},
			},
		},
		{
			name: "nested loop (not supported)",
			assembly: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"InnerLoop"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"y", "InnerLoop"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"OuterLoop"}}},
				},
				Labels: map[string]int{
					"OuterLoop": 1,
					"InnerLoop": 3,
				},
			},
			expected: [][]int{ // ループが少なくとも一つ検出できていればよい
				{2},
				{1, 2, 3},
			},
		},
		// 他のテストケースを追加...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := loop_expander.BuildControlFlowGraph(tc.assembly)
			if err != nil {
				loop_expander.PrintCFG(cfg, tc.assembly)
				t.Fatalf("BuildControlFlowGraph() error = %v", err)
			}
			loops := loop_expander.DetectLoops(cfg)
			if !reflect.DeepEqual(loops, tc.expected) {
				dot := loop_expander.ToDOT(cfg) // DOT言語を取得
				fmt.Println(dot)                // DOT言語を出力
				t.Errorf("Unexpected loops: got %v, want %v", loops, tc.expected)
			}
		})
	}
}
