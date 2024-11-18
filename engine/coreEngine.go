package engine

import (
	"fmt"
	"strconv"
	"strings"
)

// 読み込んだファイルは[]byte型なので改行区切りで[]string型に変更、コメントの削除
func filterAndFormatCode(data []byte) []string {
	lines := strings.Split(string(data), "\n")
	var result []string
	for _, line := range lines {
		if idx := strings.Index(line, "%"); idx != -1 {
			line = line[:idx]
		}

		// 前後の空白を削除
		line = strings.TrimSpace(line)

		// 空行でなければ結果に追加
		if line != "" {
			result = append(result, line)
		}
	}
	return result
}

func Execute(code []byte, speculativeWindow int) [][]string {
	traces := make([][]string, 0)
	system := systemConstructor()
	statements := filterAndFormatCode(code)
	eval(*system, statements, &traces)
	fmt.Println(traces)
	return traces
}

func eval(system System, statements []string, traces *[][]string) {
	index, _ := strconv.Atoi(system.assignment["pc"])
	if strings.Contains(statements[index], "<-") {
		// 代入文
		newSystem := assign(statements[index], system)
		eval(newSystem, statements, traces)
	} else if strings.HasPrefix(statements[index], "beqz") {
		// 条件分岐
		fmt.Println("条件分岐です")
		eval(system, statements, traces)
	} else if strings.HasPrefix(statements[index], "load") {
		// load命令
		fmt.Println("loadです")
		eval(system, statements, traces)
	} else if strings.HasPrefix(statements[index], "End:") {
		// END命令
		*traces = append(*traces, system.trace)
		fmt.Println("おわりです")
	}
}

func assign(expression string, system System) System {
	returnSystem := system
	// "x<-v<y" を "x" と "v<y" に分割
	parts := strings.SplitN(expression, "<-", 2)
	if len(parts) != 2 {
		fmt.Println("構文の解析に失敗しました: ", expression)
		return system
	}

	// 左辺と右辺をトリム
	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	// 新しくsystemを作って代入する
	returnSystem.assignment[key] = value
	return returnSystem
}

func processIf(expression string, system System, condition bool) System {
	returnSystem := system
	// スペースで分割して要素を抽出
	parts := strings.Fields(expression)
	if len(parts) != 2 || parts[0] != "beqz" {
		fmt.Println("無効な形式の条件分岐文です: ", expression)
		return returnSystem
	}

	// 変数とラベル部分をさらに分割
	subParts := strings.SplitN(parts[1], ",", 2)
	if len(subParts) != 2 {
		fmt.Println("無効な形式の条件分岐文です: ", expression)
		return returnSystem
	}

	// トリムして余計な空白を除去
	variable := strings.TrimSpace(subParts[0])
	label := strings.TrimSpace(subParts[1])

	if condition {
		returnSystem.trace = append(returnSystem.trace, "symPc(not%d)", variable)
		returnSystem.assignment["pc"]=system.assignment["pc"]+1
	} else {
		returnSystem.trace = append(returnSystem.trace, "symPc(%d)", variable)
		returnSystem.assignment["pc"]=label
	}
	return returnSystem
}

type System struct {
	memory     map[string]string
	assignment map[string]string
	trace      []string
	state      []state
}

func systemConstructor() *System {
	system := &System{
		memory:     make(map[string]string),
		assignment: make(map[string]string),
		trace:      []string{},
		state:      []state{},
	}

	// 初期値を assignment に追加
	system.assignment["pc"] = "0"

	return system
}

type state struct {
	memory     map[string]string
	assignment map[string]string
	id         int
	l          string
}
