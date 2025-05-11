package assembler

import (
    "fmt"
    "strings"
)

func DumpBasic(asm *Assembler) {
    if asm == nil {
        fmt.Println("Assembler is nil")
        return
    }

    fmt.Println("Labels:")
    for name, addr := range asm.Labels {
        fmt.Printf("  %s: %d\n", name, addr)
    }

    fmt.Println("\nProgram:")
    for _, inst := range asm.Program {
        operands := strings.Join(inst.OpCode.Operands, ", ")
        fmt.Printf("%d:\t%s %s\n", inst.Addr, inst.OpCode.Mnemonic, operands)
    }
}

func DumpBasicString(asm *Assembler) string {
	if asm == nil {
		return "Assembler is nil"
	}

	var sb strings.Builder

	sb.WriteString("Labels:\n")
	for name, addr := range asm.Labels {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", name, addr))
	}

	sb.WriteString("\nProgram:\n")
	for _, inst := range asm.Program {
		operands := strings.Join(inst.OpCode.Operands, ", ")
		sb.WriteString(fmt.Sprintf("%d:\t%s %s\n", inst.Addr, inst.OpCode.Mnemonic, operands))
	}

	return sb.String()
}

// FormatAsm はAssembler構造体からアセンブリファイルのテキスト表現を生成します。
// GenerateAsmと似ていますが、テストの比較表示に適した形式で出力します。
func FormatAsm(asm *Assembler) string {
	if asm == nil {
		return "Assembler is nil\n"
	}

	var sb strings.Builder

	// アドレスからラベル名への逆引きマップを作成 (同じアドレスに複数のラベルを格納できるようにする)
	labelsByAddr := make(map[int][]string)
	for name, addr := range asm.Labels {
		labelsByAddr[addr] = append(labelsByAddr[addr], name)
	}

	// 出力済みのラベルを記録するセット
	emittedLabels := make(map[string]bool)

	// プログラム内の命令とラベルを出力
	for _, instruction := range asm.Program {
		// アドレスに対応するラベルがあればすべて出力
		if labelNames, ok := labelsByAddr[instruction.Addr]; ok {
			for _, labelName := range labelNames {
				if !emittedLabels[labelName] {
					sb.WriteString(labelName)
					sb.WriteString(":\n")
					emittedLabels[labelName] = true
				}
			}
		}

		// 命令の出力 (インデント付き)
		sb.WriteString("\t")
		if instruction.OpCode.Mnemonic == "<-" && len(instruction.OpCode.Operands) == 2 {
			sb.WriteString(instruction.OpCode.Operands[0])
			sb.WriteString(" <- ")
			sb.WriteString(instruction.OpCode.Operands[1])
		} else {
			sb.WriteString(instruction.OpCode.Mnemonic)
			if len(instruction.OpCode.Operands) > 0 {
				operands := strings.Join(instruction.OpCode.Operands, ", ")
				if operands != "" { // spbarrのようなオペランドが空文字列の場合に対応
					sb.WriteString(" ")
					sb.WriteString(operands)
				}
			}
		}
		sb.WriteString("\n")
	}

	// プログラム末尾のラベルを出力 (命令に関連付けられていないラベル)
	// プログラムが空の場合や、最後の命令の後にラベルがある場合に対応
	maxAddr := -1
	if len(asm.Program) > 0 {
		maxAddr = asm.Program[len(asm.Program)-1].Addr
	}

	// プログラムの最大アドレスよりも大きいアドレスを持つラベルを出力
	// （ソートはしないので、ラベルの順序は保証されないが、テスト目的では十分）
	for addr, labelNames := range labelsByAddr {
		if addr > maxAddr { // 最後の命令のアドレスより大きい、またはプログラムが空(-1)の場合
			// プログラムが空でラベルのみ存在する場合も考慮
			isInstructionAtAddr := false
			for _, instr := range asm.Program {
				if instr.Addr == addr {
					isInstructionAtAddr = true
					break
				}
			}
			if !isInstructionAtAddr { // このアドレスに命令がなければラベルを出力
				for _, labelName := range labelNames {
					if !emittedLabels[labelName] {
						sb.WriteString(labelName)
						sb.WriteString(":\n")
						emittedLabels[labelName] = true
					}
				}
			}
		}
	}
	return sb.String()
}

func DumpFormatted(asm *Assembler) {
    if asm == nil {
        fmt.Println("Assembler is nil")
        return
    }

    fmt.Println("Program:")
    labelPositions := make(map[int][]string)
    for name, addr := range asm.Labels {
        labelPositions[addr] = append(labelPositions[addr], name)
    }

    for _, inst := range asm.Program {
        if labels, ok := labelPositions[inst.Addr]; ok {
            for _, label := range labels {
                fmt.Printf("%s:\n", label)
            }
        }
        operands := strings.Join(inst.OpCode.Operands, ", ")
        fmt.Printf("\t%d:\t%s %s\n", inst.Addr, inst.OpCode.Mnemonic, operands)
    }
}
