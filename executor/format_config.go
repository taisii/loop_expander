package executor

import (
	"fmt"
	"strings"
)

// PrintConfiguration 詳細なフォーマットでConfigurationを表示
func PrintConfiguration(config Configuration) {
	fmt.Println("Configuration Details:")
	fmt.Printf("  Program Counter (PC): %d\n", config.PC)
	fmt.Printf("  Step Count: %d\n", config.StepCount)

	fmt.Println("  Registers:")
	if len(config.Registers) == 0 {
		fmt.Println("    (empty)")
	} else {
		printMapStringInterface(config.Registers, "    ")
	}

	fmt.Println("  Memory:")
	if len(config.Memory) == 0 {
		fmt.Println("    (empty)")
	} else {
		printMapIntInterface(config.Memory, "    ")
	}

	PrintTrace(config.Trace)
}

func FormatConfigDifferences(expected, actual Configuration) string {
	var sb strings.Builder
	sb.WriteString("Differences between expected and actual configurations:\n")

	// PCの違い
	if expected.PC != actual.PC {
		sb.WriteString(fmt.Sprintf("- Program Counter (PC) mismatch: expected %d, got %d\n", expected.PC, actual.PC))
	}

	// StepCountの違い
	if expected.StepCount != actual.StepCount {
		sb.WriteString(fmt.Sprintf("- Step count mismatch: expected %d, got %d\n", expected.StepCount, actual.StepCount))
	}

	// Registersの違い
	if len(expected.Registers) != len(actual.Registers) {
		sb.WriteString(fmt.Sprintf("- Register count mismatch: expected %d, got %d\n",
			len(expected.Registers), len(actual.Registers)))
	} else {
		for reg, expVal := range expected.Registers {
			actVal, exists := actual.Registers[reg]
			if !exists {
				sb.WriteString(fmt.Sprintf("- Missing register in actual: %s\n", reg))
				continue
			}
			if !CompareSymbolicExpr(expVal, actVal) {
				sb.WriteString(fmt.Sprintf("- Mismatch in register %s:\n", reg))
				sb.WriteString(fmt.Sprintf("  Expected: %+v\n", expVal))
				sb.WriteString(fmt.Sprintf("  Actual:   %+v\n", actVal))
			}
		}
	}

	// Memoryの違い
	if len(expected.Memory) != len(actual.Memory) {
		sb.WriteString(fmt.Sprintf("- Memory size mismatch: expected %d, got %d\n",
			len(expected.Memory), len(actual.Memory)))
	} else {
		for addr, expVal := range expected.Memory {
			actVal, exists := actual.Memory[addr]
			if !exists {
				sb.WriteString(fmt.Sprintf("- Missing memory address in actual: %d\n", addr))
				continue
			}
			if !CompareSymbolicExpr(expVal, actVal) {
				sb.WriteString(fmt.Sprintf("- Mismatch at memory address %d:\n", addr))
				sb.WriteString(fmt.Sprintf("  Expected: %+v\n", expVal))
				sb.WriteString(fmt.Sprintf("  Actual:   %+v\n", actVal))
			}
		}
	}

	// Traceの比較
	sb.WriteString("Trace Differences:\n")
	sb.WriteString(FormatTraceDifferences(expected.Trace, actual.Trace))

	return sb.String()
}
