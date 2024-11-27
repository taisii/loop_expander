package main

import (
	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

func main() {
	program := []assembler.OpCode{
		{Mnemonic: "beqz", Operands: []string{"r1", "3"}},      // 条件: r1 == 0 (最初の分岐)
		{Mnemonic: "add", Operands: []string{"r2", "r2", "1"}}, // 偽側: r2 = r2 + 1
		{Mnemonic: "jmp", Operands: []string{"7"}},             // 偽側の終了
		{Mnemonic: "beqz", Operands: []string{"r2", "6"}},      // 真側: 条件: r2 == 0
		{Mnemonic: "add", Operands: []string{"r3", "r3", "1"}}, // 真側の偽側: r3 = r3 + 1
		{Mnemonic: "jmp", Operands: []string{"7"}},             // 真側の偽側の終了
		{Mnemonic: "add", Operands: []string{"r4", "r4", "1"}}, // 真側の真側: r4 = r4 + 1
	}
	initialConf := &executor.Configuration{
		Registers: map[string]interface{}{
			"r3": 0,
			"r4": 0,
		},
		PC:     0,
		Memory: map[int]interface{}{},
		Trace:  executor.Trace{},
	}

	finalConfigs, _ := executor.SpecExecute(program, initialConf, 100, 10)

	for _, finalConfig := range finalConfigs {
		executor.PrintConfiguration(*finalConfig)
	}
}
