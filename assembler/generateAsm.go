package assembler

import "strings"

// GenerateAsm はAssembler構造体からアセンブリファイルのテキスト表現を生成します。
func GenerateAsm(assembler *Assembler) (string, error) {
	var sb strings.Builder

	// アドレスからラベル名への逆引きマップを作成 (同じアドレスに複数のラベルを格納できるようにする)
	labelsByAddr := make(map[int][]string)
	for name, addr := range assembler.Labels {
		labelsByAddr[addr] = append(labelsByAddr[addr], name)
	}

	// 出力済みのラベルを記録するセット
	emittedLabels := make(map[string]bool)

	for _, instruction := range assembler.Program {
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

		// 命令の出力
		if instruction.OpCode.Mnemonic == "<-" && len(instruction.OpCode.Operands) == 2 {
			sb.WriteString(instruction.OpCode.Operands[0])
			sb.WriteString(" <- ")
			sb.WriteString(instruction.OpCode.Operands[1])
		} else {
			sb.WriteString(instruction.OpCode.Mnemonic)
			if len(instruction.OpCode.Operands) > 0 {
				sb.WriteString(" ")
				sb.WriteString(strings.Join(instruction.OpCode.Operands, ", "))
			}
		}
		sb.WriteString("\n")
	}

	// プログラム末尾のラベルを出力 (命令に関連付けられていないラベル)
	if len(assembler.Program) > 0 {
		lastInstructionAddr := assembler.Program[len(assembler.Program)-1].Addr + 1 // 最後の命令の次のアドレス
		for addr, labelNames := range labelsByAddr {
			if addr >= lastInstructionAddr {
				for _, labelName := range labelNames {
					if !emittedLabels[labelName] {
						sb.WriteString(labelName)
						sb.WriteString(":\n")
						emittedLabels[labelName] = true
					}
				}
			}
		}
	} else {
		// プログラムが空の場合でもラベルを出力する
		for _, labelNames := range labelsByAddr {
			for _, labelName := range labelNames {
				if !emittedLabels[labelName] {
					sb.WriteString(labelName)
					sb.WriteString(":\n")
					emittedLabels[labelName] = true
				}
			}
		}
	}

	return sb.String(), nil
}
