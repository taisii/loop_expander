package assembler

import (
	"fmt"
	"strconv"
)

// 新しいアセンブラを作成
func NewAssembler() *Assembler {
	return &Assembler{
		Program: []Program{},
		Labels:  make(map[string]int),
	}
}

// プログラムをロードし、命令にアドレスを割り当て、ラベルを解決
func (a *Assembler) LoadProgram(ins []Instruction) {
	a.Program = []Program{}
	a.Labels = make(map[string]int)
	addr := 0

	for _, inst := range ins {
		switch v := inst.(type) {
		case Label:
			// ラベルの場合、アドレスを保存
			a.Labels[v.Name] = addr
		default:
			// 命令の場合、プログラムに追加し、アドレスをインクリメント
			a.Program = append(a.Program, Program{Addr: addr, Inst: inst})
			addr++
		}
	}

	// オペランド中のラベルをアドレスに解決
	for i, program := range a.Program {
		if op, ok := program.Inst.(OpCode); ok {
			for j, operand := range op.Operands {
				// オペランドが数値か判定
				if _, err := strconv.Atoi(operand); err != nil {
					// 数値でない場合
					if labelAddr, exists := a.Labels[operand]; exists {
						// ラベルとして解決
						op.Operands[j] = strconv.Itoa(labelAddr)
					} else {
						// 必要なら追加の検証をここに入れられる
					}
				}
			}
			a.Program[i].Inst = op
		}
	}
}

// プログラムを表示
func (a *Assembler) ShowProgram() {
	for _, prog := range a.Program {
		fmt.Printf("  %d: %s\n", prog.Addr, prog.Inst)
	}
	for name, addr := range a.Labels {
		fmt.Printf("sym %s = %d\n", name, addr)
	}
}
