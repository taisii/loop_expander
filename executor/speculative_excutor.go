package executor

import (
	"github.com/taisii/go-project/assembler"
)

func SpeculativeExecution(
	conf *Configuration,
	program []assembler.OpCode,
	maxSteps, maxSpecDepth int,
) ([]*Configuration, error) {
	executionState := &ExecutionState{
		Counter:     0,
		CurrentConf: *conf, // 初期状態を設定
		Speculative: []SpeculativeState{},
	}

	completedConfs := []*Configuration{} // 探索が完了した構成を格納
	stepCount := 0

	for stepCount < maxSteps {
		// 現在の構成を取得
		currntConf := executionState.CurrentConf

		if currntConf.PC < 0 || currntConf.PC >= len(program) {
			// プログラムカウンタが範囲外の場合は完了状態として保存
			completedConfs = append(completedConfs, &currntConf)
			break
		}

		// 現在の命令を取得
		inst := program[currntConf.PC]

		// 投機状態の処理
		if len(executionState.Speculative) > 0 {
			topSpec := &executionState.Speculative[0]

			if topSpec.RemainingWin == 0 {
				// 投機的ウィンドウが終了 -> コミットまたはロールバック
				executionState.Speculative = executionState.Speculative[1:] // スタックから削除

				if currntConf.PC == topSpec.CorrectPC {
					// コミット
					currntConf.Trace.Observations = append(currntConf.Trace.Observations, Observation{
						Type:      ObsTypeCommit,
						PC:        topSpec.CorrectPC,
						SpecState: topSpec,
					})
				} else {
					// ロールバック
					currntConf = topSpec.InitialConf // 初期状態にロールバック
					currntConf.Trace.Observations = append(currntConf.Trace.Observations, Observation{
						Type:      ObsTypeRollback,
						PC:        topSpec.CorrectPC,
						SpecState: topSpec,
					})
				}
				// 更新後の構成を現在の実行状態にセット
				executionState.CurrentConf = currntConf
				continue
			}

			// 投機的ウィンドウをデクリメント
			topSpec.RemainingWin--
		}

		// 現在の命令を実行
		newConfs, newSpecStates, err := AlwaysMispredictStep(inst, executionState, maxSpecDepth)
		if err != nil {
			return nil, err
		}

		// 投機的状態を追加
		for _, newSpec := range newSpecStates {
			executionState.Speculative = append([]SpeculativeState{newSpec}, executionState.Speculative...)
		}

		// 新しい構成を反映
		if len(newConfs) > 0 {
			executionState.CurrentConf = *newConfs[0] // 最初の構成を現在の状態に設定
			newConfs = newConfs[1:]                   // 残りをキューに保持
		}

		// 完了した構成を保存
		completedConfs = append(completedConfs, newConfs...)

		stepCount++
	}

	return completedConfs, nil
}
