package executor

import (
	"fmt"

	"github.com/taisii/go-project/assembler"
)

// AlwaysMispredictStep handles a single instruction under always-mispredict semantics.
func AlwaysMispredictStep(
	inst assembler.OpCode,
	curConf *Configuration,
) ([]*Configuration, bool, error) {
	newConfs := []*Configuration{}
	isSpeculative := false // 投機実行が必要かどうかを示すフラグ

	switch inst.Mnemonic {
	case "beqz":
		// beqz reg, target
		if len(inst.Operands) != 2 {
			return nil, false, fmt.Errorf("beqz requires 2 operands, got %d", len(inst.Operands))
		}
		reg := inst.Operands[0]
		target, err := evalExpr(inst.Operands[1], curConf)
		if err != nil {
			return nil, false, err
		}
		condition, err := evalExpr(reg, curConf)
		if err != nil {
			return nil, false, err
		}

		// トレースイベントの初期化
		traceEventTrue := Observation{
			Type: ObsTypePC,
			Value: SymbolicExpr{
				Op:       "==",
				Operands: []interface{}{reg, 0},
			},
			PC: curConf.PC,
		}
		traceEventFalse := Observation{
			Type: ObsTypePC,
			Value: SymbolicExpr{
				Op:       "!=",
				Operands: []interface{}{reg, 0},
			},
			PC: curConf.PC,
		}

		switch condValue := condition.(type) {
		case int:
			// Concrete condition
			newConf := *curConf
			if condValue == 0 {
				// Condition true
				isSpeculative = true
				newConf.PC++
				newConf.Trace.PathCond = updatePathCond(newConf.Trace.PathCond, "==", reg)
				newConf.Trace.Observations = append(newConf.Trace.Observations, traceEventFalse) // 誤ってfalseのほうに進むからトレースはfalseのもの
			} else {
				// Condition false,
				isSpeculative = true
				newConf.PC = int(target.(int))
				newConf.Trace.PathCond = updatePathCond(newConf.Trace.PathCond, "!=", reg)
				newConf.Trace.Observations = append(newConf.Trace.Observations, traceEventTrue) // 誤ってtrueのほうに進むからトレースはtrueのもの
			}
			newConfs = append(newConfs, &newConf)

		case SymbolicExpr:
			// Symbolic condition
			isSpeculative = true
			newConfTrue := *curConf
			newConfFalse := *curConf

			// True branch (condition is true, mispredicts to False)
			newConfTrue.PC++
			newConfTrue.Trace.PathCond = updatePathCond(curConf.Trace.PathCond, "==", reg)
			newConfTrue.Trace.Observations = append(newConfTrue.Trace.Observations, traceEventFalse) // 誤ってfalseのほうに進むからトレースはfalseのもの

			// False branch (condition is false, mispredicts to True)
			newConfFalse.PC = int(target.(int))
			newConfFalse.Trace.PathCond = updatePathCond(curConf.Trace.PathCond, "!=", reg)
			newConfFalse.Trace.Observations = append(newConfFalse.Trace.Observations, traceEventTrue) // 誤ってtrueのほうに進むからトレースはtrueのもの

			// 両方の分岐を返す
			newConfs = append(newConfs, &newConfTrue, &newConfFalse)

		default:
			return nil, false, fmt.Errorf("unexpected type for condition: %T", condValue)
		}

	default:
		// Unsupported instructions are handled with the default step
		conf, err := Step(inst, curConf)
		if err != nil {
			return nil, false, err
		}
		newConfs = append(newConfs, conf...)
	}

	return newConfs, isSpeculative, nil
}
