package main

import (
	"github.com/taisii/go-project/assembler" // パッケージをインポート
)

func main() {
	// 例としてμAsmの命令を定義
	// ins := []Instruction{
	// 	Label{Name: "start"},
	// 	OpCode{Mnemonic: "mov", Operands: []string{"r1", "r2"}},
	// 	OpCode{Mnemonic: "add", Operands: []string{"r1", "r3"}},
	// 	Label{Name: "loop"},
	// 	OpCode{Mnemonic: "jmp", Operands: []string{"start"}},
	// }

	ins2 := []assembler.Instruction{
		assembler.Label{Name: "start"},
		assembler.OpCode{Mnemonic: "cmp", Operands: []string{"x", "v", "y"}}, // x <- v < y
		assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "End"}},   // branch if x == 0 to End
		assembler.OpCode{Mnemonic: "load", Operands: []string{"v", "v"}},     // load value from array1
		assembler.OpCode{Mnemonic: "load", Operands: []string{"v", "v"}},     // load value from array2
		assembler.Label{Name: "End"},
	}

	// アセンブラを作成し、プログラムをロード
	assembler := assembler.NewAssembler()
	assembler.LoadProgram(ins2)

	// プログラムを表示
	assembler.ShowProgram()
}
