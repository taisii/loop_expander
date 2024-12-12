package loop_expander

import "github.com/taisii/go-project/assembler"

// BasicBlock は、基本ブロックを表す構造体
type BasicBlock struct {
	StartAddress int                     // 開始アドレス
	EndAddress   int                     // 終了アドレス
	Instructions []assembler.Instruction // 命令のリスト
	Succs        []int                   // 後続ブロックの開始アドレス
}

// ControlFlowGraph は、制御フローグラフを表す構造体
type ControlFlowGraph struct {
	Blocks []*BasicBlock // 基本ブロックのリスト
}
