package loop_expander

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/taisii/go-project/assembler"
)

// Loop_expander関数
func Loop_expander(asm *assembler.Assembler, maxUnrollCount int) (*assembler.Assembler, error) {
	if asm == nil || maxUnrollCount <= 0 {
		return nil, errors.New("invalid arguments")
	}

	cfg, err := BuildControlFlowGraph(asm)
	if err != nil {
		return nil, fmt.Errorf("failed to build CFG: %w", err)
	}

	loops := DetectLoops(cfg)

	if len(loops) == 0 {
		return asm, nil
	}

	StartAddress := loops[0][0]
	loopEndAddress := loops[0][len(loops[0])-1]
	loopProgram := asm.Program[StartAddress:]
	loopLength := len(loopProgram)

	expandedAsm := assembler.CopyAssembler(asm)
	expandedAsm.Program = expandedAsm.Program[:StartAddress]

	for i := 0; i < maxUnrollCount; i++ {
		for _, inst := range loopProgram {
			var nextAddr int
			if len(expandedAsm.Program) > 0 {
				nextAddr = expandedAsm.Program[len(expandedAsm.Program)-1].Addr + 1
			} else {
				nextAddr = 0 // プログラムの先頭のアドレスを設定
			}

			newInst := assembler.Instruction{
				Addr: nextAddr,
				OpCode: assembler.OpCode{
					Mnemonic: inst.OpCode.Mnemonic,
					Operands: make([]string, len(inst.OpCode.Operands)),
				},
			}
			copy(newInst.OpCode.Operands, inst.OpCode.Operands)

			// 新しいラベル名とアドレスを生成()
			newLabels := make(map[string]int)
			for label, addr := range asm.Labels {
				newLabelName := label
				newLabelAddr := addr
				if addr >= StartAddress && addr <= loopEndAddress {
					newLabelName = fmt.Sprintf("%s_%d", label, i)
					newLabelAddr = addr - StartAddress + (loopLength * (i + 1))
					newLabels[newLabelName] = newLabelAddr
				}
			}

			// 元のラベルを新しいラベルに置き換え
			for j, operand := range inst.OpCode.Operands {
				for originalLabel := range asm.Labels {
					if operand == originalLabel {
						if _, ok := newLabels[operand+"_"+strconv.Itoa(i)]; ok {
							newInst.OpCode.Operands[j] = operand + "_" + strconv.Itoa(i)
						}
					}
				}
			}
			// 新しいラベルをexpandedAsmに追加
			for newLabelName, newLabelAddr := range newLabels {
				expandedAsm.Labels[newLabelName] = newLabelAddr + StartAddress
			}
			expandedAsm.Program = append(expandedAsm.Program, newInst)

		}
	}

	return expandedAsm, nil
}
