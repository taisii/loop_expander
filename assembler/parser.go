package assembler

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// ParseAsm はμAsmのアセンブリのファイルを読み込み、Assembler構造体に変換します。
func ParseAsm(r io.Reader) (*Assembler, error) {
	assembler := &Assembler{
		Program: make([]Instruction, 0),
		Labels:  make(map[string]int),
	}
	scanner := bufio.NewScanner(r)
	addr := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "%") { // 空行またはコメント行はスキップ
			continue
		}
		if strings.HasSuffix(line, ":") { // ラベル行
			labelName := strings.TrimSuffix(line, ":")
			assembler.Labels[labelName] = addr
		} else { // 命令行
			// 代入命令の判定と分割
			if strings.Contains(line, "<-") {
				parts := strings.Split(line, "<-")
				if len(parts) == 2 {
					mnemonic := "<-"
					operands := []string{strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])}
					assembler.Program = append(assembler.Program, Instruction{Addr: addr, OpCode: OpCode{Mnemonic: mnemonic, Operands: operands}})
					addr++
					continue
				}
			}

			// spbarr
			if strings.Contains(line, "spbarr") {
				fmt.Println("spbarr", line)
				mnemonic := "spbarr"
				assembler.Program = append(assembler.Program, Instruction{Addr: addr, OpCode: OpCode{
					Mnemonic: mnemonic,
				}})
				addr++
				continue
			}

			// 命令行
			parts := strings.Fields(line)
			if len(parts) > 1 { // ニーモニックとオペランドがあることを確認
				mnemonic := parts[0]
				operands := []string{parts[1]} // 初期値として最初の要素をオペランドに設定

				// カンマが含まれていればカンマで分割
				if strings.Contains(operands[0], ",") {
					operands = strings.Split(operands[0], ",")
					// オペランドの前後の空白を削除
					for i := range operands {
						operands[i] = strings.TrimSpace(operands[i])
					}
				}

				assembler.Program = append(assembler.Program, Instruction{Addr: addr, OpCode: OpCode{Mnemonic: mnemonic, Operands: operands}})
				addr++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("アセンブリファイルのスキャン中にエラーが発生しました: %w", err)
	}
	return assembler, nil
}
