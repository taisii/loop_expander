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