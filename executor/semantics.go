package executor

import (
	"errors"
	"fmt"

	"github.com/taisii/go-project/assembler"
)

// Step executes a single instruction
func Step(inst assembler.OpCode, conf *Configuration) ([]*Configuration, error) {
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

		// Check if condition is concrete (int) or symbolic
		switch condValue := condition.(type) {
		case int:
			if condValue == 0 {
				// True branch
				conf.PC = int(target.(int))
				return []*Configuration{conf}, nil
			} else {
				// False branch
				conf.PC++
				return []*Configuration{conf}, nil
			}
		case SymbolicExpr:
			// Generate two configurations: one for true, one for false
			confTrue := *conf
			confFalse := *conf

			confTrue.Registers = copyRegisters(conf.Registers)
			confTrue.PC = int(target.(int))
			confFalse.Registers = copyRegisters(conf.Registers)
			confFalse.PC++

			return []*Configuration{
				&confTrue,
				&confFalse,
			}, nil
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

func ExecuteProgram(program []assembler.OpCode, configuration *Configuration, maxSteps int) error {
	queue := []*Configuration{configuration}
	steps := 0 // 現在のステップ数を追跡

	for len(queue) > 0 {
		if steps >= maxSteps {
			return errors.New("maximum steps reached")
		}

		current := queue[0]
		queue = queue[1:]

		if current.PC >= len(program) {
			continue // プログラム終了
		}

		inst := program[current.PC]
		newConfigs, err := Step(inst, current)
		if err != nil {
			return err
		}

		queue = append(queue, newConfigs...)
		steps++ // ステップ数をインクリメント
	}

	return nil
}
