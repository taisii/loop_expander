package loop_expander

import (
	"fmt"
	"sort"
	"strings"

	"github.com/taisii/go-project/assembler"
)

// BuildControlFlowGraph は、アセンブリプログラムから制御フローグラフを構築
func BuildControlFlowGraph(asm *assembler.Assembler) (*ControlFlowGraph, error) {
	// 1. 基本ブロックへの分割
	blocks, err := splitToBasicBlocks(asm)
	if err != nil {
		return nil, err
	}

	// 2. 制御フローグラフの構築
	cfg := &ControlFlowGraph{
		Blocks: blocks,
	}
	buildCFGEdges(cfg, asm)

	return cfg, nil
}

// splitToBasicBlocks は、アセンブリプログラムを基本ブロックに分割
func splitToBasicBlocks(assembler *assembler.Assembler) ([]*BasicBlock, error) {
	blocks := make([]*BasicBlock, 0)
	block := &BasicBlock{}
	for _, instruction := range assembler.Program {
		// ラベルの直後または新しいブロックの開始時にブロックを開始
		isLabel := false
		for _, labelAddr := range assembler.Labels {
			if instruction.Addr == labelAddr {
				isLabel = true
				break
			}
		}
		if isLabel || len(block.Instructions) == 0 {
			if len(block.Instructions) > 0 {
				blocks = append(blocks, block)
			}
			block = &BasicBlock{StartAddress: instruction.Addr}
		}

		block.Instructions = append(block.Instructions, instruction)
		block.EndAddress = instruction.Addr

		// ジャンプ命令または分岐命令の場合、新しいブロックを開始
		if instruction.OpCode.Mnemonic == "jmp" || instruction.OpCode.Mnemonic == "beqz" {
			blocks = append(blocks, block)
			block = &BasicBlock{StartAddress: instruction.Addr + 1}
		}
	}
	if len(block.Instructions) > 0 { // 最後のブロックを追加
		blocks = append(blocks, block)
	}
	return blocks, nil
}

// buildCFGEdges は、制御フローグラフのエッジを構築します。
func buildCFGEdges(cfg *ControlFlowGraph, asm *assembler.Assembler) {
	for i, block := range cfg.Blocks {
		// 最後の命令がジャンプ命令または分岐命令の場合
		if len(block.Instructions) > 0 {
			lastInst := block.Instructions[len(block.Instructions)-1]
			if lastInst.OpCode.Mnemonic == "jmp" || lastInst.OpCode.Mnemonic == "beqz" {
				labelName := lastInst.OpCode.Operands[len(lastInst.OpCode.Operands)-1]
				labelAddr, ok := asm.Labels[labelName]
				if ok {
					// ジャンプ先のブロック番号を後続ブロックに追加
					blockIndex := findBlockIndexByAddr(cfg, labelAddr)
					if blockIndex != -1 {
						block.Succs = append(block.Succs, blockIndex)
					} else {
						// ラベルに対応するブロックが見つからない場合はエラー処理
						// Succsを空にする
						block.Succs = []int{}
					}
				} else {
					// ラベルが見つからない場合はエラー処理
					fmt.Printf("ラベル %s が見つかりません\n", labelName)
				}
			}
			if lastInst.OpCode.Mnemonic != "jmp" && i < len(cfg.Blocks)-1 {
				// 次のブロック番号を後続ブロックに追加
				block.Succs = append(block.Succs, findBlockIndexByAddr(cfg, cfg.Blocks[i+1].StartAddress))
			}
		}

		// Succs リストをソート
		sort.Ints(block.Succs)
	}
}

// findBlockIndexByAddr は、指定されたアドレスを持つブロックのインデックスを返します。
func findBlockIndexByAddr(cfg *ControlFlowGraph, addr int) int {
	for i, block := range cfg.Blocks {
		if block.StartAddress == addr {
			return i
		}
	}
	return -1 // 見つからない場合は -1 を返す
}

// PrintCFG は、ControlFlowGraph を見やすい形で出力する関数です。
func PrintCFG(cfg *ControlFlowGraph, asm *assembler.Assembler) {
	for i, block := range cfg.Blocks {
		fmt.Printf("Block %d (Addr: %d-%d):\n", i, block.StartAddress, block.EndAddress)
		for _, inst := range block.Instructions {
			fmt.Printf("  %s\n", inst.String())
		}
		fmt.Printf("  Succs: ")
		for _, succAddr := range block.Succs {
			fmt.Printf("%d ", succAddr)
			// 後続ブロックがラベルの場合、ラベル名を表示
			for labelName, labelAddr := range asm.Labels {
				if labelAddr == succAddr {
					fmt.Printf("(%s) ", labelName)
					break
				}
			}
		}
		fmt.Println()
	}
}

func ToDOT(cfg *ControlFlowGraph) string {
	var sb strings.Builder
	sb.WriteString("digraph CFG {\n")
	for i, block := range cfg.Blocks {
		// ノードのラベルを生成 (ブロック番号とアドレスのみ)
		nodeLabel := fmt.Sprintf("Block %d\nAddr: %d-%d", i, block.StartAddress, block.EndAddress)

		// ノードを記述
		sb.WriteString(fmt.Sprintf("  %d [label=\"%s\"];\n", i, nodeLabel))

		// エッジを記述
		for _, succ := range block.Succs {
			sb.WriteString(fmt.Sprintf("  %d -> %d;\n", i, succ))
		}
	}
	sb.WriteString("}\n")
	return sb.String()
}
