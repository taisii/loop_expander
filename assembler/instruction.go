package assembler

import (
	"fmt"
	"strings"
)

// ラベルを表す構造体
type Label struct {
	Name string
	Addr int
}

// 実際の命令を表す構造体
type OpCode struct {
	Mnemonic string
	Operands []string
}

// p/2に対応：アドレスと命令のペアを保持
type Instruction struct {
	Addr   int
	OpCode OpCode
}

// μAsmアセンブラを表す構造体
type Assembler struct {
	Program []Instruction  // 命令とアドレスのペアのリスト
	Labels  map[string]int // ラベル名とアドレスのマップ
}

func (inst Instruction) String() string {
	return fmt.Sprintf("Addr: %d, OpCode: %s", inst.Addr, inst.OpCode.String())
}

// Instructionインターフェースの実装：OpCode
func (op OpCode) String() string {
	return fmt.Sprintf("%s %s", op.Mnemonic, strings.Join(op.Operands, ", "))
}

// Instructionインターフェースの実装：Label
func (l Label) String() string {
	return fmt.Sprintf("label %s", l.Name)
}
