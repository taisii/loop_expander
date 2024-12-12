package assembler

import "fmt"

// CompareInstructions は、2つの Instruction を比較する関数です。
func CompareInstructions(a, b Instruction) bool {
	return a.Addr == b.Addr && a.OpCode.String() == b.OpCode.String()
}

// diffInstructions は、2つの Instruction を比較し、差分をわかりやすく出力する関数です。
func DiffInstructions(a, b Instruction) string {
	if a.Addr != b.Addr {
		return fmt.Sprintf("Addr differs: got %d, want %d", a.Addr, b.Addr)
	}

	// OpCode のフィールドを再帰的に比較
	return DiffOpCodes(a.OpCode, b.OpCode)
}

// diffOpCodes は、2つの OpCode を比較し、差分をわかりやすく出力する関数です。
func DiffOpCodes(a, b OpCode) string {
	if a.Mnemonic != b.Mnemonic {
		return fmt.Sprintf("Mnemonic differs: got %s, want %s", a.Mnemonic, b.Mnemonic)
	}
	if len(a.Operands) != len(b.Operands) {
		return fmt.Sprintf("Number of operands differs: got %d, want %d", len(a.Operands), len(b.Operands))
	}
	for i, operandA := range a.Operands {
		if operandA != b.Operands[i] {
			return fmt.Sprintf("Operand %d differs: got %s, want %s", i, operandA, b.Operands[i])
		}
	}
	return "" // 差分なし
}
