package executor

import (
	"fmt"

	"github.com/taisii/go-project/assembler"
)

// AlwaysMispredictStep handles a single instruction under always-mispredict semantics.
func AlwaysMispredictStep(
	inst assembler.OpCode,
	execState *ExecutionState,
	maxSpecDepth int,
) ([]*Configuration, []SpeculativeState, error) {
	curConf := &execState.CurrentConf
	newConfs := []*Configuration{}
	newSpecStates := []SpeculativeState{}

	switch inst.Mnemonic {
	case "beqz":
		// beqz reg, target
		if len(inst.Operands) != 2 {
			return nil, nil, fmt.Errorf("beqz requires 2 operands, got %d", len(inst.Operands))
		}
		reg := inst.Operands[0]
		target, err := evalExpr(inst.Operands[1], curConf)
		if err != nil {
			return nil, nil, err
		}
		condition, err := evalExpr(reg, curConf)
		if err != nil {
			return nil, nil, err
		}

		// トレースイベントの初期化
		traceEvent := Observation{
			Type: ObsTypePC,
			Value: SymbolicExpr{
				Op:       "==",
				Operands: []interface{}{reg, 0},
			},
			PC: curConf.PC,
		}

		switch condValue := condition.(type) {
		case int:
			// Concrete condition
			if condValue == 0 {
				// 誤予測として False branch を進む
				curConf.PC++
				curConf.Trace.PathCond = updatePathCond(curConf.Trace.PathCond, "!=", reg)
			} else {
				// 誤予測として True branch を進む
				curConf.PC = int(target.(int))
				curConf.Trace.PathCond = updatePathCond(curConf.Trace.PathCond, "==", reg)
			}

			curConf.Trace.Observations = append(curConf.Trace.Observations, traceEvent)
			newConfs = append(newConfs, curConf)

		case SymbolicExpr:
			// Symbolic condition
			if len(execState.Speculative) >= maxSpecDepth {
				// 最大深度に達している場合は現在の分岐を進む
				confCopy := *curConf
				confCopy.PC++
				confCopy.Trace.PathCond = updatePathCond(curConf.Trace.PathCond, "==", reg)
				confCopy.Trace.Observations = append(confCopy.Trace.Observations, traceEvent)
				newConfs = append(newConfs, &confCopy)

				return newConfs, newSpecStates, nil
			}

			// 投機的実行のための新しい状態を準備
			confFalse := *curConf
			confTrue := *curConf

			// False branch (condition is true, mispredicts to False)
			confFalse.PC++ // 次の命令へ進む
			confFalse.Trace.PathCond = updatePathCond(curConf.Trace.PathCond, "==", reg)
			confFalse.Trace.Observations = append(confFalse.Trace.Observations, Observation{
				Type:      ObsTypePC,
				PC:        confFalse.PC,
				Value:     reg,
				SpecState: nil,
			})
			specStateFalse := SpeculativeState{
				ID:           execState.Counter,
				RemainingWin: maxSpecDepth,
				StartPC:      curConf.PC,
				InitialConf:  confFalse,
				CorrectPC:    confTrue.PC, // 誤予測のためTrueへ進むべき
			}

			// True branch (condition is false, mispredicts to True)
			confTrue.PC = int(target.(int)) // ジャンプ先へ進む
			confTrue.Trace.PathCond = updatePathCond(curConf.Trace.PathCond, "!=", reg)
			confTrue.Trace.Observations = append(confTrue.Trace.Observations, Observation{
				Type:      ObsTypePC,
				PC:        confTrue.PC,
				Value:     reg,
				SpecState: nil,
			})
			specStateTrue := SpeculativeState{
				ID:           execState.Counter + 1,
				RemainingWin: maxSpecDepth,
				StartPC:      curConf.PC,
				InitialConf:  confTrue,
				CorrectPC:    confFalse.PC, // 誤予測のためFalseへ進むべき
			}

			// 投機的状態を保存
			newSpecStates = append(newSpecStates, specStateFalse, specStateTrue)
			execState.Counter += 2 // IDを2つ消費

			// 誤予測のFalse branchを現在の構成として進める
			newConfs = append(newConfs, &confFalse)
			return newConfs, newSpecStates, nil

		default:
			return nil, nil, fmt.Errorf("unexpected type for condition: %T", condValue)
		}

	case "jmp":
		// jmp target
		if len(inst.Operands) != 1 {
			return nil, nil, fmt.Errorf("jmp requires 1 operand, got %d", len(inst.Operands))
		}
		target, err := evalExpr(inst.Operands[0], curConf)
		if err != nil {
			return nil, nil, err
		}
		curConf.PC = int(target.(int))

		curConf.Trace.Observations = append(curConf.Trace.Observations, Observation{
			Type: ObsTypePC,
			PC:   curConf.PC,
		})
		newConfs = append(newConfs, curConf)

	default:
		// Use default step for unsupported instructions
		conf, err := Step(inst, curConf)
		if err != nil {
			return nil, nil, err
		}
		newConfs = conf
	}

	return newConfs, newSpecStates, nil
}
