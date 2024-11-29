package executor

import (
	"fmt"
)

// 初期状態と最終状態を受け取ってそれぞれの状態とトレースを出力
func PrintTest(initialConfig, finalConfig Configuration) {
	// Assignments 出力
	fmt.Println("Assignments:")
	printMapStringInterface(finalConfig.Registers, "  ")
	fmt.Println()

	// 初期状態の出力
	fmt.Println("initial conf:")
	printMemoryAndRegister(initialConfig)

	// トレースの出力
	PrintTrace(finalConfig.Trace)

	// 最終状態の出力
	fmt.Println("\nfinal conf:")
	printMemoryAndRegister(finalConfig)

	// Path Condition の出力
	fmt.Println("\nPath Condition:")
	fmt.Printf("  %s\n", formatSymbolicExpr(finalConfig.Trace.PathCond))
	fmt.Println("===========================")
}
