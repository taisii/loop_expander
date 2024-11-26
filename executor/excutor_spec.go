package executor

import (
	"errors"

	"github.com/taisii/go-project/assembler"
)

// 実行パスの定義
type ExecutionPath struct {
	CurrentConf      Configuration      // Current configuration
	SpeculativeStack []SpeculativeState // Stack of speculative states
}

func initializePaths(initialConfig *Configuration) []ExecutionPath {
	if initialConfig.Registers == nil {
		initialConfig.Registers = make(map[string]interface{})
	}
	if initialConfig.Memory == nil {
		initialConfig.Memory = make(map[int]interface{})
	}
	return []ExecutionPath{
		{
			CurrentConf:      *initialConfig,
			SpeculativeStack: nil,
		},
	}
}

func handleRollback(currentConf Configuration, specState SpeculativeState) Configuration {
	rollbackConf := Configuration{
		PC:        specState.CorrectPC,
		Registers: copyRegisters(specState.Configuration.Registers),
		Memory:    copyMemory(specState.Configuration.Memory),
		Trace:     currentConf.Trace,
	}

	// ロールバック操作をトレースに追加
	rollbackConf.Trace.Observations = append(
		rollbackConf.Trace.Observations,
		Observation{
			Type:  ObsTypeRollback,
			Value: specState.ID,
			PC:    rollbackConf.PC,
		},
	)

	return rollbackConf
}

func handleSpecStart(newConfs []*Configuration, correctConfs []*Configuration, path ExecutionPath, defaultRemainingWindow int) []ExecutionPath {
	var newPaths []ExecutionPath

	for i, conf := range newConfs {
		newSpecState := SpeculativeState{
			ID:            len(path.SpeculativeStack),
			RemainingWin:  defaultRemainingWindow,
			StartPC:       path.CurrentConf.PC,
			Configuration: path.CurrentConf,
			CorrectPC:     correctConfs[i].PC, // 正しい実行と投機的実行の条件探索順序が同じであることを仮定
		}

		// SpeculativeStack が空でない場合は RemainingWin を計算
		if len(path.SpeculativeStack) > 0 {
			newSpecState.RemainingWin = path.SpeculativeStack[len(path.SpeculativeStack)-1].RemainingWin - 1
		}

		newPath := ExecutionPath{
			CurrentConf:      *conf,
			SpeculativeStack: append(path.SpeculativeStack, newSpecState),
		}
		observation := Observation{
			PC:    path.CurrentConf.PC,
			Type:  ObsTypeStart,
			Value: newSpecState.ID,
		}

		// Observations に新しい要素を最後から2番目に挿入
		obs := newPath.CurrentConf.Trace.Observations
		if len(obs) >= 1 {
			newPath.CurrentConf.Trace.Observations = append(obs[:len(obs)-1], observation, obs[len(obs)-1])
		} else {
			newPath.CurrentConf.Trace.Observations = append(obs, observation)
		}

		newPaths = append(newPaths, newPath)
	}

	return newPaths
}

// execute runs the given program with the provided initial configuration up to maxSteps.
func SpecExecute(program []assembler.OpCode, initialConfig *Configuration, maxSteps int, remainingWindow int) ([]*Configuration, error) {

	paths := initializePaths(initialConfig)
	var finalConfigs []*Configuration
	stepCount := 0

	// 実行ループ
	for stepCount < maxSteps {
		stepCount++

		if len(paths) > 0 {
			// 現在の状態を確認
			currentPath := paths[len(paths)-1]
			paths = paths[:len(paths)-1]

			// プログラム終了判定
			if currentPath.CurrentConf.PC >= len(program) {
				if len(currentPath.SpeculativeStack) > 0 {
					// ロールバック処理
					lastSpecState := currentPath.SpeculativeStack[len(currentPath.SpeculativeStack)-1]
					currentPath.SpeculativeStack = currentPath.SpeculativeStack[:len(currentPath.SpeculativeStack)-1]
					currentPath.CurrentConf = handleRollback(currentPath.CurrentConf, lastSpecState)

					paths = append(paths, currentPath)
				} else {
					// 実行完了
					finalConfigs = append(finalConfigs, &currentPath.CurrentConf)
				}
				continue
			}

			// 命令実行フェーズ
			instruction := program[currentPath.CurrentConf.PC]
			newConfs, isSpeculative, err := AlwaysMispredictStep(instruction, &currentPath.CurrentConf)
			if err != nil {
				return nil, err
			}

			if isSpeculative {
				//ここでStep関数を実行して正しい遷移先を取得している。2つのsemanticsを表す関数が同じ順序でconfsを返すことが前提になっている
				correctConfs, err := Step(instruction, &currentPath.CurrentConf)
				if err != nil {
					return nil, err
				}
				paths = append(paths, handleSpecStart(newConfs, correctConfs, currentPath, remainingWindow)...)
			} else {
				// 通常の命令実行
				currentPath.CurrentConf = *newConfs[0]
				paths = append(paths, currentPath)
			}
		} else {
			return finalConfigs, nil
		}

		// 最大ステップ数を超えた場合のエラー処理
		if stepCount >= maxSteps {
			return nil, errors.New("execution reached maximum step limit")
		}
	}

	return finalConfigs, nil
}
