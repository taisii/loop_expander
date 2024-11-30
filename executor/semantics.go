package executor

import (
	"fmt"

	"github.com/taisii/go-project/assembler"
)

// Step executes a single instruction
func Step(instruction assembler.OpCode, conf *Configuration) ([]*Configuration, error) {
	copiedConf := copyConfiguration(*conf)
	var traceEvent Observation    // トレースイベントを初期化
	traceEvent.PC = copiedConf.PC // 現在のプログラムカウンタを設定

	switch instruction.Mnemonic {
	case "mov":
		// mov dest, src
		if len(instruction.Operands) != 2 {
			return nil, fmt.Errorf("mov requires 2 operands, got %d", len(instruction.Operands))
		}
		dest := instruction.Operands[0]
		srcValue, err := evalExpr(instruction.Operands[1], &copiedConf)
		if err != nil {
			return nil, err
		}
		copiedConf.Registers[dest] = srcValue
		copiedConf.PC++

		// トレースイベントを追加
		traceEvent.Type = ObsTypeStore
		traceEvent.Address = &SymbolicExpr{Op: "var", Operands: []interface{}{dest}}
		traceEvent.Value = srcValue
		copiedConf.Trace.Observations = append(copiedConf.Trace.Observations, traceEvent)

		return []*Configuration{&copiedConf}, nil

	case "add":
		// add dest, src1, src2
		if len(instruction.Operands) != 3 {
			return nil, fmt.Errorf("add requires 3 operands, got %d", len(instruction.Operands))
		}
		dest := instruction.Operands[0]
		src1, err := evalExpr(instruction.Operands[1], &copiedConf)
		if err != nil {
			return nil, err
		}
		src2, err := evalExpr(instruction.Operands[2], &copiedConf)
		if err != nil {
			return nil, err
		}
		result, err := evalExpr(SymbolicExpr{
			Op:       "+",
			Operands: []interface{}{src1, src2},
		}, &copiedConf)
		if err != nil {
			return nil, err
		}
		copiedConf.Registers[dest] = result
		copiedConf.PC++

		// トレースイベントを追加
		traceEvent.Type = ObsTypeStore
		traceEvent.Address = &SymbolicExpr{Op: "var", Operands: []interface{}{dest}}
		traceEvent.Value = result
		copiedConf.Trace.Observations = append(copiedConf.Trace.Observations, traceEvent)

		return []*Configuration{&copiedConf}, nil

	case "beqz":
		// beqz reg, target
		if len(instruction.Operands) != 2 {
			return nil, fmt.Errorf("beqz requires 2 operands, got %d", len(instruction.Operands))
		}
		target, err := evalExpr(instruction.Operands[1], &copiedConf)
		if err != nil {
			return nil, err
		}
		reg, err := evalExpr(instruction.Operands[0], &copiedConf)
		if err != nil {
			return nil, err
		}

		traceEventTrue := Observation{
			PC:   conf.PC,
			Type: ObsTypePC,
			Value: SymbolicExpr{
				Op:       "==",
				Operands: []interface{}{reg, 0},
			},
		}
		traceEventFalse := Observation{
			PC:   conf.PC,
			Type: ObsTypePC,
			Value: SymbolicExpr{
				Op:       "!=",
				Operands: []interface{}{reg, 0},
			},
		}

		switch condValue := reg.(type) {
		case int:
			// Concrete condition
			if condValue == 0 {
				// True branch
				copiedConf.PC = int(target.(int))
				copiedConf.Trace.PathCond = updatePathCond(copiedConf.Trace.PathCond, "==", reg)
				copiedConf.Trace.Observations = append(copiedConf.Trace.Observations, traceEventTrue)
				return []*Configuration{&copiedConf}, nil
			} else {
				// False branch
				copiedConf.PC++
				copiedConf.Trace.PathCond = updatePathCond(copiedConf.Trace.PathCond, "!=", reg)
				copiedConf.Trace.Observations = append(copiedConf.Trace.Observations, traceEventFalse)
				return []*Configuration{&copiedConf}, nil
			}
		case SymbolicExpr:
			// Symbolic condition
			confTrue := copiedConf
			confFalse := copiedConf

			// Copy registers and PathCond for each branch
			confTrue.Registers = copyRegisters(copiedConf.Registers)
			confFalse.Registers = copyRegisters(copiedConf.Registers)

			// True branch
			confTrue.PC = int(target.(int))
			confTrue.Trace.PathCond = updatePathCond(copiedConf.Trace.PathCond, "==", reg)
			confTrue.Trace.Observations = append(confTrue.Trace.Observations, traceEventTrue)

			// False branch
			confFalse.PC++
			confFalse.Trace.PathCond = updatePathCond(copiedConf.Trace.PathCond, "!=", reg)
			confFalse.Trace.Observations = append(confFalse.Trace.Observations, traceEventFalse)

			return []*Configuration{&confTrue, &confFalse}, nil

		default:
			return nil, fmt.Errorf("unexpected type for condition: %T", condValue)
		}

	case "load":
		// load dest, addr
		if len(instruction.Operands) != 2 {
			return nil, fmt.Errorf("load requires 2 operands, got %d", len(instruction.Operands))
		}
		dest := instruction.Operands[0]
		addrExpr, err := ParseSymbolicExpr(instruction.Operands[1])
		if err != nil {
			return nil, err
		}

		// アドレス式を評価
		addrValue, err := evalExpr(*addrExpr, &copiedConf)
		if err != nil {
			return nil, err
		}

		// メモリから値を取得
		address, ok := addrValue.(int)
		if !ok {
			return nil, fmt.Errorf("address must be an integer, got %T", addrValue)
		}
		value, exists := copiedConf.Memory[address]
		if !exists {
			return nil, fmt.Errorf("memory address %d not found", address)
		}

		// 値をレジスタに保存
		copiedConf.Registers[dest] = value
		copiedConf.PC++

		// トレースイベントを追加
		traceEvent.Type = ObsTypeLoad
		traceEvent.Address = address
		traceEvent.Value = value
		copiedConf.Trace.Observations = append(copiedConf.Trace.Observations, traceEvent)

		return []*Configuration{&copiedConf}, nil

	case "store":
		// store value, addr
		if len(instruction.Operands) != 2 {
			return nil, fmt.Errorf("store requires 2 operands, got %d", len(instruction.Operands))
		}
		valueExpr := instruction.Operands[0]
		addressExpression, err := ParseSymbolicExpr(instruction.Operands[1])
		if err != nil {
			return nil, err
		}

		// 値とアドレスを評価
		value, err := evalExpr(valueExpr, &copiedConf)
		if err != nil {
			return nil, err
		}
		addrValue, err := evalExpr(*addressExpression, &copiedConf)
		if err != nil {
			return nil, err
		}

		address, ok := addrValue.(int)
		if !ok {
			return nil, fmt.Errorf("address must be an integer, got %T", addrValue)
		}

		// メモリを更新
		copiedConf.Memory[address] = value
		copiedConf.PC++

		// トレースイベントを追加
		traceEvent.Type = ObsTypeStore
		traceEvent.Address = addrValue
		traceEvent.Value = value
		copiedConf.Trace.Observations = append(copiedConf.Trace.Observations, traceEvent)

		return []*Configuration{&copiedConf}, nil

	case "jmp":
		// jmp target
		if len(instruction.Operands) != 1 {
			return nil, fmt.Errorf("jmp requires 1 operand, got %d", len(instruction.Operands))
		}
		target, err := evalExpr(instruction.Operands[0], &copiedConf)
		if err != nil {
			return nil, err
		}
		copiedConf.PC = int(target.(int))

		// トレースイベントを追加
		traceEvent.Type = ObsTypePC
		traceEvent.Value = SymbolicExpr{Op: "jmp", Operands: []interface{}{target}}
		copiedConf.Trace.Observations = append(copiedConf.Trace.Observations, traceEvent)

		return []*Configuration{&copiedConf}, nil

	default:
		return nil, fmt.Errorf("unsupported instruction: %s", instruction.Mnemonic)
	}
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
