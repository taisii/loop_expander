package executor

import (
	"github.com/taisii/go-project/assembler"
)

func ExecuteProgram(program []assembler.OpCode, configuration *Configuration, maxSteps int) ([]*Configuration, error) {
	// キューに初期状態を追加（各パスごとに個別のステップカウントを保持）
	queue := []*Configuration{configuration}
	completedConfigs := []*Configuration{} // 完了したすべての状態を収集

	for len(queue) > 0 {
		// キューから現在の状態を取得
		current := queue[0]
		queue = queue[1:]

		// パスのステップ数が最大値を超えた場合、このパスを破棄
		if current.StepCount >= maxSteps {
			continue
		}

		// プログラム終了時に最終状態を収集
		if current.PC >= len(program) {
			completedConfigs = append(completedConfigs, current)
			continue
		}

		// 現在の命令を取得
		inst := program[current.PC]

		// 命令を実行し、新しい状態を取得
		newConfigs, err := Step(inst, current)
		if err != nil {
			return nil, err
		}

		// 新しい状態に対してステップカウントをインクリメントし、キューに追加
		for _, newConfig := range newConfigs {
			newConfig.StepCount = current.StepCount + 1 // 現在のステップ数を引き継ぎ＋1
			queue = append(queue, newConfig)
		}
	}

	return completedConfigs, nil
}
