package assembler

import (
	"errors"
	"fmt"
	"strconv"
)

// Configuration structure
type Configuration struct {
	Memory    map[string]int
	Registers map[string]int
	PC        int // Program Counter
}

// NewConfiguration creates a new Configuration
func NewConfiguration(memory map[string]int, registers map[string]int) *Configuration {
	return &Configuration{
		Memory:    memory,
		Registers: registers,
		PC:        0,
	}
}

// Instruction Evaluation
func evalExpr(expr interface{}, conf *Configuration) (int, error) {
	switch v := expr.(type) {
	case int: // Direct value
		return v, nil
	case string: // Register, PC, or a parsable integer
		if v == "pc" { // PCの場合
			return conf.PC, nil
		}
		if val, ok := conf.Registers[v]; ok { // レジスタの場合
			return val, nil
		}
		if intValue, err := strconv.Atoi(v); err == nil { // 整数にパース可能ならそれを返す
			return intValue, nil
		}
		return 0, fmt.Errorf("unknown register or symbol: %s", v)
	case []interface{}: // Operation (e.g., ["+", "x", 1])
		if len(v) < 2 {
			return 0, errors.New("invalid operation")
		}
		op, args := v[0].(string), v[1:]
		switch op {
		case "+": // 加算
			x, err := evalExpr(args[0], conf)
			if err != nil {
				return 0, err
			}
			y, err := evalExpr(args[1], conf)
			if err != nil {
				return 0, err
			}
			return x + y, nil
		case "-": // 減算
			x, err := evalExpr(args[0], conf)
			if err != nil {
				return 0, err
			}
			y, err := evalExpr(args[1], conf)
			if err != nil {
				return 0, err
			}
			return x - y, nil
		default:
			return 0, fmt.Errorf("unsupported operation: %s", op)
		}
	default:
		return 0, fmt.Errorf("unsupported expression type: %T", v)
	}
}

// Step executes a single instruction
func Step(inst OpCode, conf *Configuration) (*Configuration, error) {
	switch inst.Mnemonic {
	case "mov":
		if len(inst.Operands) != 2 {
			return nil, fmt.Errorf("mov requires 2 operands, got %d", len(inst.Operands))
		}
		dest := inst.Operands[0]
		value, err := evalExpr(inst.Operands[1], conf)
		if err != nil {
			return nil, err
		}
		conf.Registers[dest] = value
		conf.PC++
	case "add":
		if len(inst.Operands) != 3 {
			return nil, fmt.Errorf("add requires 3 operands, got %d", len(inst.Operands))
		}
		dest := inst.Operands[0]
		x, err := evalExpr(inst.Operands[1], conf)
		if err != nil {
			return nil, err
		}
		y, err := evalExpr(inst.Operands[2], conf)
		if err != nil {
			return nil, err
		}
		conf.Registers[dest] = x + y
		conf.PC++
	case "beqz":
		if len(inst.Operands) != 2 {
			return nil, fmt.Errorf("beqz requires 2 operands, got %d", len(inst.Operands))
		}
		reg := inst.Operands[0]
		target, err := evalExpr(inst.Operands[1], conf)
		if err != nil {
			return nil, err
		}
		if conf.Registers[reg] == 0 {
			conf.PC = target
		} else {
			conf.PC++
		}
	case "jmp":
		if len(inst.Operands) != 1 {
			return nil, fmt.Errorf("jmp requires 1 operand, got %d", len(inst.Operands))
		}
		target, err := evalExpr(inst.Operands[0], conf)
		if err != nil {
			return nil, err
		}
		conf.PC = target
	default:
		return nil, fmt.Errorf("unsupported instruction: %s", inst.Mnemonic)
	}
	return conf, nil
}

// Run executes the program until completion or timeout
func Run(program []OpCode, conf *Configuration, maxSteps int) (*Configuration, error) {
	for steps := 0; steps < maxSteps; steps++ {
		if conf.PC < 0 || conf.PC >= len(program) {
			return conf, nil // Program finished
		}
		inst := program[conf.PC]
		var err error
		conf, err = Step(inst, conf)
		if err != nil {
			return nil, err
		}
	}
	return conf, errors.New("timeout reached")
}
