package executor

import (
	"fmt"

	"github.com/taisii/go-project/assembler"
)

// Step executes a single instruction
func Step(inst assembler.OpCode, conf *Configuration) ([]*Configuration, error) {
	var traceEvent Observation // トレースイベントを初期化
	traceEvent.PC = conf.PC    // 現在のプログラムカウンタを設定

	switch inst.Mnemonic {
	case "mov":
		// mov dest, src
		if len(inst.Operands) != 2 {
			return nil, fmt.Errorf("mov requires 2 operands, got %d", len(inst.Operands))
		}
		dest := inst.Operands[0]
		srcValue, err := evalExpr(inst.Operands[1], conf)
		if err != nil {
			return nil, err
		}
		conf.Registers[dest] = srcValue
		conf.PC++

		// トレースイベントを追加
		traceEvent.Type = ObsTypeStore
		traceEvent.Address = &SymbolicExpr{Op: "var", Operands: []interface{}{dest}}
		traceEvent.Value = srcValue
		conf.Trace.Observations = append(conf.Trace.Observations, traceEvent)

		return []*Configuration{conf}, nil

	case "add":
		// add dest, src1, src2
		if len(inst.Operands) != 3 {
			return nil, fmt.Errorf("add requires 3 operands, got %d", len(inst.Operands))
		}
		dest := inst.Operands[0]
		src1, err := evalExpr(inst.Operands[1], conf)
		if err != nil {
			return nil, err
		}
		src2, err := evalExpr(inst.Operands[2], conf)
		if err != nil {
			return nil, err
		}
		result, err := evalExpr(SymbolicExpr{
			Op:       "+",
			Operands: []interface{}{src1, src2},
		}, conf)
		if err != nil {
			return nil, err
		}
		conf.Registers[dest] = result
		conf.PC++

		// トレースイベントを追加
		traceEvent.Type = ObsTypeStore
		traceEvent.Address = &SymbolicExpr{Op: "var", Operands: []interface{}{dest}}
		traceEvent.Value = result
		conf.Trace.Observations = append(conf.Trace.Observations, traceEvent)

		return []*Configuration{conf}, nil

	case "beqz":
		// beqz reg, target
		if len(inst.Operands) != 2 {
			return nil, fmt.Errorf("beqz requires 2 operands, got %d", len(inst.Operands))
		}
		reg := inst.Operands[0]
		target, err := evalExpr(inst.Operands[1], conf)
		if err != nil {
			return nil, err
		}
		condition, err := evalExpr(reg, conf)
		if err != nil {
			return nil, err
		}

		traceEvent.Type = ObsTypePC
		traceEvent.Value = SymbolicExpr{
			Op:       "==",
			Operands: []interface{}{reg, 0},
		}

		switch condValue := condition.(type) {
		case int:
			// Concrete condition
			if condValue == 0 {
				// True branch
				conf.PC = int(target.(int))
				conf.Trace.PathCond = SymbolicExpr{
					Op:       "&&",
					Operands: []interface{}{conf.Trace.PathCond, SymbolicExpr{Op: "==", Operands: []interface{}{reg, 0}}},
				}
				conf.Trace.Observations = append(conf.Trace.Observations, traceEvent)
				return []*Configuration{conf}, nil
			} else {
				// False branch
				conf.PC++
				conf.Trace.PathCond = SymbolicExpr{
					Op:       "&&",
					Operands: []interface{}{conf.Trace.PathCond, SymbolicExpr{Op: "!=", Operands: []interface{}{reg, 0}}},
				}
				conf.Trace.Observations = append(conf.Trace.Observations, traceEvent)
				return []*Configuration{conf}, nil
			}
		case SymbolicExpr:
			// Symbolic condition
			confTrue := *conf
			confFalse := *conf

			// Copy registers and PathCond for each branch
			confTrue.Registers = copyRegisters(conf.Registers)
			confFalse.Registers = copyRegisters(conf.Registers)

			// True branch
			confTrue.PC = int(target.(int))
			confTrue.Trace.PathCond = updatePathCond(conf.Trace.PathCond, "==", reg)
			confTrue.Trace.Observations = append(confTrue.Trace.Observations, traceEvent)

			// False branch
			confFalse.PC++
			confFalse.Trace.PathCond = updatePathCond(conf.Trace.PathCond, "!=", reg)
			confFalse.Trace.Observations = append(confFalse.Trace.Observations, traceEvent)

			return []*Configuration{&confTrue, &confFalse}, nil

		default:
			return nil, fmt.Errorf("unexpected type for condition: %T", condValue)
		}

	case "jmp":
		// jmp target
		if len(inst.Operands) != 1 {
			return nil, fmt.Errorf("jmp requires 1 operand, got %d", len(inst.Operands))
		}
		target, err := evalExpr(inst.Operands[0], conf)
		if err != nil {
			return nil, err
		}
		conf.PC = int(target.(int))

		// トレースイベントを追加
		traceEvent.Type = ObsTypePC
		traceEvent.Value = SymbolicExpr{Op: "jmp", Operands: []interface{}{target}}
		conf.Trace.Observations = append(conf.Trace.Observations, traceEvent)

		return []*Configuration{conf}, nil

	default:
		return nil, fmt.Errorf("unsupported instruction: %s", inst.Mnemonic)
	}
}

func copyRegisters(registers map[string]interface{}) map[string]interface{} {
	newRegisters := make(map[string]interface{})
	for k, v := range registers {
		newRegisters[k] = v
	}
	return newRegisters
}

func updatePathCond(currentCond SymbolicExpr, op string, reg interface{}) SymbolicExpr {
	newCond := SymbolicExpr{
		Op:       op,
		Operands: []interface{}{reg, 0},
	}

	if currentCond.Op == "" && len(currentCond.Operands) == 0 {
		// 現在の条件が空の場合、新しい条件をそのまま返す
		return newCond
	}

	// 現在の条件がある場合、新しい条件と連結
	return SymbolicExpr{
		Op:       "&&",
		Operands: []interface{}{currentCond, newCond},
	}
}
