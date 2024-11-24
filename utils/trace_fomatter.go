package utils

import (
	"fmt"
	"strings"

	"github.com/taisii/go-project/executor"
)

// PrintTrace 詳細なフォーマットでTraceを表示
func PrintTrace(trace executor.Trace) {
	fmt.Println("Trace :")

	// 観測データの表示
	if len(trace.Observations) == 0 {
		fmt.Println("  Observations: (none)")
	} else {
		fmt.Println("  Observations:")
		for _, obs := range trace.Observations {
			printObservation(obs) // 既存のprintObservationを利用
		}
	}

	// パス条件の表示
	fmt.Println("  Path Condition:")
	if len(trace.PathCond.Operands) == 0 {
		fmt.Println("    (none)")
	} else {
		fmt.Printf("    %s\n", formatSymbolicExpr(trace.PathCond))
	}
	fmt.Println("===========================")
}

func FormatTraceDifferences(expected, actual executor.Trace) string {
	var sb strings.Builder
	sb.WriteString("Differences between expected and actual traces:\n")

	// 観測数の違い
	if len(expected.Observations) != len(actual.Observations) {
		sb.WriteString(fmt.Sprintf("- Observation count mismatch: expected %d, got %d\n",
			len(expected.Observations), len(actual.Observations)))
	} else {
		// 各観測の比較
		for i := 0; i < len(expected.Observations); i++ {
			expectedObs := expected.Observations[i]
			actualObs := actual.Observations[i]

			if expectedObs.PC != actualObs.PC {
				sb.WriteString(fmt.Sprintf("- Mismatch at observation %d (PC): expected %d, got %d\n",
					i+1, expectedObs.PC, actualObs.PC))
			}
			if expectedObs.Type != actualObs.Type {
				sb.WriteString(fmt.Sprintf("- Mismatch at observation %d (Type): expected %s, got %s\n",
					i+1, expectedObs.Type, actualObs.Type))
			}
			if !executor.CompareSymbolicExpr(expectedObs.Address, actualObs.Address) {
				sb.WriteString(fmt.Sprintf("- Mismatch at observation %d (Address):\n", i+1))
				sb.WriteString(fmt.Sprintf("  Expected: %+v\n", expectedObs.Address))
				sb.WriteString(fmt.Sprintf("  Actual:   %+v\n", actualObs.Address))
			}
			if !executor.CompareSymbolicExpr(expectedObs.Value, actualObs.Value) {
				sb.WriteString(fmt.Sprintf("- Mismatch at observation %d (Value):\n", i+1))
				sb.WriteString(fmt.Sprintf("  Expected: %+v\n", expectedObs.Value))
				sb.WriteString(fmt.Sprintf("  Actual:   %+v\n", actualObs.Value))
			}
		}
	}

	// パス条件の違い
	if !executor.CompareSymbolicExpr(expected.PathCond, actual.PathCond) {
		sb.WriteString("- Path condition mismatch:\n")
		sb.WriteString(fmt.Sprintf("  Expected: %+v\n", expected.PathCond))
		sb.WriteString(fmt.Sprintf("  Actual:   %+v\n", actual.PathCond))
	}

	return sb.String()
}

// printObservation 観測データを整形して出力
func printObservation(obs executor.Observation) {
	fmt.Printf("  PC: %d, Type: %s", obs.PC, obs.Type)

	// Addressがある場合の処理
	if obs.Address != nil {
		fmt.Printf(", Address: %s", formatValue(obs.Address))
	}

	// Valueがある場合の処理
	if obs.Value != nil {
		fmt.Printf(", Value: %s", formatValue(obs.Value))
	}

	// SpeculativeStateがある場合の処理
	if obs.SpecState != nil {
		fmt.Printf(", SpeculativeState: {ID: %d, RemainingWin: %d, StartPC: %d, CorrectPC: %d, InitialConf: {Registers: %v, Memory: %v}}",
			obs.SpecState.ID,
			obs.SpecState.RemainingWin,
			obs.SpecState.StartPC,
			obs.SpecState.CorrectPC,
			obs.SpecState.InitialConf.Registers, // Assume InitialConf contains Registers and Memory
			obs.SpecState.InitialConf.Memory)
	}

	fmt.Println()
}

// printMemoryAndRegister Configuration を整形して出力
func printMemoryAndRegister(config executor.Configuration) {
	fmt.Println("  m=")
	printMapStringInterface(config.Registers, "    ")
	fmt.Println("  a=")
	printMapIntInterface(config.Memory, "    ")
}

// printMapStringInterface map[string]interface{} を整形して出力
func printMapStringInterface(data map[string]interface{}, indent string) {
	for key, value := range data {
		fmt.Printf("%s%s: %s\n", indent, key, formatValue(value))
	}
}

// printMapIntInterface map[int]interface{} を整形して出力
func printMapIntInterface(data map[int]interface{}, indent string) {
	for key, value := range data {
		fmt.Printf("%s%d: %s\n", indent, key, formatValue(value))
	}
}
