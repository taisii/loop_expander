package executor

import (
	"errors"
	"fmt"

	"github.com/taisii/go-project/assembler"
)

// execute runs the given program with the provided initial configuration up to maxSteps.
func SpecExecute(program []assembler.OpCode, initialConfig *Configuration, maxSteps int) ([]*Configuration, error) {
	// 初期化
	execState := ExecutionState{
		Counter:     0,
		CurrentConf: *initialConfig,
		Speculative: nil, // 初期状態では投機状態なし
	}
	var finalConfigs []*Configuration
	stepCount := 0

	// 実行ループ
	for stepCount < maxSteps {
		stepCount++

		// 現在の状態を確認
		var currentConf *Configuration
		if len(execState.Speculative) > 0 {
			// 投機的状態が存在する場合はスタックの先頭を取得
			currentSpec := &execState.Speculative[len(execState.Speculative)-1]
			currentConf = &currentSpec.InitialConf
		} else {
			// 非投機的状態を処理
			currentConf = &execState.CurrentConf
		}

		// プログラム終了条件
		if currentConf.PC >= len(program) || currentConf.PC < 0 {
			finalConfigs = append(finalConfigs, currentConf)
			if len(execState.Speculative) > 0 {
				execState.Speculative = execState.Speculative[:len(execState.Speculative)-1] // 投機的状態を削除
			} else {
				break
			}
			continue
		}

		// 現在の命令を取得
		instruction := program[currentConf.PC]

		// step関数を呼び出して命令を実行
		newConfs, isSpeculative, err := AlwaysMispredictStep(instruction, currentConf)
		if err != nil {
			return nil, fmt.Errorf("execution error at step %d: %w", stepCount, err)
		}

		if len(execState.Speculative) > 0 {
			// 投機的状態にいる場合の処理
			currentSpec := &execState.Speculative[len(execState.Speculative)-1]
			currentSpec.RemainingWin--

			if currentSpec.RemainingWin <= 0 {
				// 投機ウィンドウが終了
				if currentSpec.CorrectPC == currentConf.PC {
					// コミット
					execState.CurrentConf = currentSpec.InitialConf
					finalConfigs = append(finalConfigs, &execState.CurrentConf)
				}
				// スタックから削除
				execState.Speculative = execState.Speculative[:len(execState.Speculative)-1]
			} else {
				// 投機的状態の進行
				execState.Speculative[len(execState.Speculative)-1].InitialConf = *newConfs[0]
			}
		} else if isSpeculative {
			// 非投機的状態から投機的状態に移行
			for _, conf := range newConfs {
				specState := SpeculativeState{
					ID:           execState.Counter,
					RemainingWin: 5, // 固定の投機ウィンドウ（適宜調整可能）
					StartPC:      currentConf.PC,
					InitialConf:  *conf,
					CorrectPC:    conf.PC, // 仮に正しい分岐を現在のPCとする
				}
				execState.Speculative = append(execState.Speculative, specState)
				execState.Counter++
			}
		} else {
			// 非投機的状態を更新
			execState.CurrentConf = *newConfs[0]
		}
	}

	// 最大ステップ数を超えた場合のエラー処理
	if stepCount >= maxSteps {
		return nil, errors.New("execution reached maximum step limit")
	}

	return finalConfigs, nil
}
