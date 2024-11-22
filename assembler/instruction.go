package assembler

import (
	"fmt"
	"strings"
)

// 命令を表すインターフェース
type Instruction interface {
	String() string
}

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
type Program struct {
	Addr int
	Inst Instruction
}

// μAsmアセンブラを表す構造体
type Assembler struct {
	Program []Program      // 命令とアドレスのペアのリスト
	Labels  map[string]int // ラベル名とアドレスのマップ
}

// Instructionインターフェースの実装：OpCode
func (op OpCode) String() string {
	return fmt.Sprintf("%s %s", op.Mnemonic, strings.Join(op.Operands, ", "))
}

// Instructionインターフェースの実装：Label
func (l Label) String() string {
	return fmt.Sprintf("label %s", l.Name)
}
