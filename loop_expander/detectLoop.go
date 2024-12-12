package loop_expander

// detectLoops は、制御フローグラフに対してDFSを行い、ループを検出します。
func DetectLoops(cfg *ControlFlowGraph) [][]int {
	visited := make(map[int]bool) // 訪問済みノードを記録するマップ
	stack := make([]int, 0)       // DFSのためのスタック
	loops := make([][]int, 0)     // 検出されたループを格納するリスト
	for i := range cfg.Blocks {   // 各ブロックを始点としてDFS
		detectLoopsDFS(cfg, i, visited, stack, &loops)
	}
	return loops
}

// detectLoopsDFS は、detectLoops のための深さ優先探索を行う補助関数です。
func detectLoopsDFS(cfg *ControlFlowGraph, blockIndex int, visited map[int]bool, stack []int, loops *[][]int) {
	if visited[blockIndex] { // 既に訪問済み
		return
	}
	visited[blockIndex] = true
	stack = append(stack, blockIndex)

	for _, succIndex := range cfg.Blocks[blockIndex].Succs {
		if i := indexOf(stack, succIndex); i != -1 { // スタックに含まれている場合、ループを検出
			loop := stack[i:] // ループを構成するブロックのインデックスを抽出
			*loops = append(*loops, loop)
		} else {
			detectLoopsDFS(cfg, succIndex, visited, stack, loops)
		}
	}

	stack = stack[:len(stack)-1] // スタックから現在のブロックを削除
}

// indexOf は、スライス内で指定した要素が最初に出現するインデックスを返します。
func indexOf(s []int, e int) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}
