package loop_expander

// reachableSubGraph は、開始ノードから到達可能なノードで構成される部分グラフを返します。
func reachableSubGraph(graph *ControlFlowGraph, startIndex int) *ControlFlowGraph {
	if graph == nil || len(graph.Blocks) == 0 || startIndex < 0 || startIndex >= len(graph.Blocks) {
		return &ControlFlowGraph{} // 空のグラフを返す
	}

	reachable := make(map[int]bool)
	queue := []int{startIndex}
	reachable[startIndex] = true

	for len(queue) > 0 {
		currentIndex := queue[0]
		queue = queue[1:]

		currentBlock := graph.Blocks[currentIndex]
		for _, succAddr := range currentBlock.Succs {
			succIndex := findBlockIndexByAddress(graph, succAddr)
			if succIndex != -1 && !reachable[succIndex] {
				reachable[succIndex] = true
				queue = append(queue, succIndex)
			}
		}
	}

	subGraph := &ControlFlowGraph{Blocks: make([]*BasicBlock, 0)}
	for i := 0; i < len(graph.Blocks); i++ {
		if reachable[i] {
			subGraph.Blocks = append(subGraph.Blocks, graph.Blocks[i])
		}
	}

	return subGraph
}

// アドレスからブロックのインデックスを検索するヘルパー関数
func findBlockIndexByAddress(graph *ControlFlowGraph, address int) int {
	for i, block := range graph.Blocks {
		if block.StartAddress == address {
			return i
		}
	}
	return -1 // 見つからない場合は-1を返す
}
