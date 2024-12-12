package loop_expander

import (
	"testing"

	"github.com/taisii/go-project/assembler"
)

func TestBuildControlFlowGraph(t *testing.T) {
	testCases := []struct {
		name     string
		assembly *assembler.Assembler
		expected *ControlFlowGraph
	}{
		{
			name: "Simple program",
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
			expected: &ControlFlowGraph{
				Blocks: []*BasicBlock{
					{
						StartAddress: 0,
						EndAddress:   2,
						Instructions: []assembler.Instruction{
							{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
							{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
							{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "L1"}}},
						},
						Succs: []int{1,2},
					},
					{
						StartAddress: 3,
						EndAddress:   3,
						Instructions: []assembler.Instruction{
							{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"z", "2"}}},
						},
						Succs: []int{2},
					},
					{
						StartAddress: 4,
						EndAddress:   4,
						Instructions: []assembler.Instruction{
							{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"w", "3"}}},
						},
						Succs: []int{},
					},
				},
			},
		},
		{
			name: "Program with loop",
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
			expected: &ControlFlowGraph{
				Blocks: []*BasicBlock{
					{
						StartAddress: 0,
						EndAddress:   1,
						Instructions: []assembler.Instruction{
							{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
							{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"y", "1"}}},
						},
						Succs: []int{1},
					},
					{
						StartAddress: 2,
						EndAddress:   2,
						Instructions: []assembler.Instruction{
							{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "L1"}}},
						},
						Succs: []int{2, 3},
					},
					{
						StartAddress: 3,
						EndAddress:   4,
						Instructions: []assembler.Instruction{
							{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"z", "2"}}},
							{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"L2"}}},
						},
						Succs: []int{1},
					},
					{
						StartAddress: 5,
						EndAddress:   6,
						Instructions: []assembler.Instruction{
							{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"w", "3"}}},
							{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"v", "4"}}},
						},
						Succs: []int{},
					},
				},
			},
		},
		// 他のテストケースを追加...
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := BuildControlFlowGraph(tc.assembly)
			if err != nil {
				t.Fatalf("BuildControlFlowGraph() error = %v", err)
			}

			if len(cfg.Blocks) != len(tc.expected.Blocks) {
				// エラー発生時に CFG を出力
				PrintCFG(cfg, tc.assembly)
				t.Fatalf("Unexpected number of blocks: got %d, want %d", len(cfg.Blocks), len(tc.expected.Blocks))
			}
			for i, block := range cfg.Blocks {
				expectedBlock := tc.expected.Blocks[i]
				if block.StartAddress != expectedBlock.StartAddress {
					// エラー発生時に CFG を出力
					PrintCFG(cfg, tc.assembly)
					t.Errorf("Block %d: unexpected StartAddress: got %d, want %d", i, block.StartAddress, expectedBlock.StartAddress)
				}
				if block.EndAddress != expectedBlock.EndAddress {
					// エラー発生時に CFG を出力
					PrintCFG(cfg, tc.assembly)
					t.Errorf("Block %d: unexpected EndAddress: got %d, want %d", i, block.EndAddress, expectedBlock.EndAddress)
				}
				if len(block.Instructions) != len(expectedBlock.Instructions) {
					// エラー発生時に CFG を出力
					PrintCFG(cfg, tc.assembly)
					t.Fatalf("Block %d: unexpected number of instructions: got %d, want %d", i, len(block.Instructions), len(expectedBlock.Instructions))
				}
				for j, inst := range block.Instructions {
					if diff := assembler.DiffInstructions(inst, expectedBlock.Instructions[j]); diff != "" { // diffInstructions 関数を使用
						// エラー発生時に CFG を出力
						PrintCFG(cfg, tc.assembly)
						t.Errorf("Block %d, instruction %d: %s", i, j, diff)
					}
				}
				if len(block.Succs) != len(expectedBlock.Succs) {
					// エラー発生時に CFG を出力
					PrintCFG(cfg, tc.assembly)
					t.Fatalf("Block %d: unexpected number of successors: got %d, want %d", i, len(block.Succs), len(expectedBlock.Succs))
				}
				for j, succ := range block.Succs {
					if succ != expectedBlock.Succs[j] {
						// エラー発生時に CFG を出力
						PrintCFG(cfg, tc.assembly)
						t.Errorf("Block %d, successor %d: unexpected successor: got %d, want %d", i, j, succ, expectedBlock.Succs[j])
					}
				}
			}
		})
	}
}
