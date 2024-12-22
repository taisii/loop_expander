package assembler

import (
	"fmt"
	"reflect"
	"strings"
)

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

// CompareAssembler は、2つの Assembler を比較する関数です。
func CompareAssembler(a, b *Assembler) bool {
	if a == nil || b == nil {
		return a == nil && b == nil // 両方nilなら等しい
	}

	if len(a.Program) != len(b.Program) || !reflect.DeepEqual(a.Labels, b.Labels) {
		return false
	}
	for i := range a.Program {
		if !CompareInstructions(a.Program[i], b.Program[i]) {
			return false
		}
	}
	return true
}

// DiffAssembler は、2つの Assembler を比較し、差分をわかりやすく出力する関数です。
func DiffAssembler(a, b *Assembler) string {
	if a == nil && b == nil {
		return "" // 両方nilなら差分なし
	}
	if a == nil {
		return "First Assembler is nil"
	}
	if b == nil {
		return "Second Assembler is nil"
	}

	if !reflect.DeepEqual(a.Labels, b.Labels) {
		var gotLabels strings.Builder
		gotLabels.WriteString("got Labels:\n")
		for label, addr := range a.Labels {
			gotLabels.WriteString(fmt.Sprintf("    %s: %d\n", label, addr))
		}
		var wantLabels strings.Builder
		wantLabels.WriteString("want Labels:\n")
		for label, addr := range b.Labels {
			wantLabels.WriteString(fmt.Sprintf("    %s: %d\n", label, addr))
		}
		return fmt.Sprintf("Label map diff:\n%s%s", gotLabels.String(), wantLabels.String()) // 改行を追加
	}

	if len(a.Program) != len(b.Program) {
		return fmt.Sprintf("Program length differs: got %d, want %d", len(a.Program), len(b.Program))
	}

	for i := range a.Program {
		if diff := DiffInstructions(a.Program[i], b.Program[i]); diff != "" {
			return fmt.Sprintf("Instruction at index %d differs: %s", i, diff)
		}
	}

	return "" // 差分なし
}
