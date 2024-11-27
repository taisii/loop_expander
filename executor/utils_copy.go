package executor

func copyRegisters(registers map[string]interface{}) map[string]interface{} {
	newRegisters := make(map[string]interface{})
	for k, v := range registers {
		newRegisters[k] = v
	}
	return newRegisters
}

func copyMemory(memory map[int]interface{}) map[int]interface{} {
	newMemory := make(map[int]interface{})
	for addr, value := range memory {
		newMemory[addr] = value
	}
	return newMemory
}

func copyConfiguration(conf Configuration) Configuration {
	newRegisters := copyRegisters(conf.Registers)
	newMemory := copyMemory(conf.Memory)

	// Trace のコピー
	newTrace := Trace{
		Observations: make([]Observation, len(conf.Trace.Observations)),
		PathCond:     copySymbolicExpr(conf.Trace.PathCond), // SymbolicExpr のディープコピー
	}
	for i, obs := range conf.Trace.Observations {
		newTrace.Observations[i] = copyObservation(obs) // Observation のディープコピー
	}

	return Configuration{
		PC:        conf.PC,
		Registers: newRegisters,
		Memory:    newMemory,
		Trace:     newTrace,
		StepCount: conf.StepCount,
	}
}

func copySymbolicExpr(expr SymbolicExpr) SymbolicExpr {
	newOperands := make([]interface{}, len(expr.Operands))
	for i, operand := range expr.Operands {
		switch v := operand.(type) {
		case SymbolicExpr:
			newOperands[i] = copySymbolicExpr(v)
		default:
			newOperands[i] = v
		}
	}

	return SymbolicExpr{
		Op:       expr.Op,
		Operands: newOperands,
	}
}

func copyObservation(obs Observation) Observation {
	var newSpecState *SpeculativeState
	if obs.SpecState != nil {
		newSpecStateCopy := copySpecState(*obs.SpecState)
		newSpecState = &newSpecStateCopy
	}

	return Observation{
		PC:        obs.PC,
		Type:      obs.Type,
		Address:   obs.Address, // `interface{}` 型だが値型の可能性が高い
		Value:     obs.Value,   // 同上
		SpecState: newSpecState,
	}
}

func copySpecState(state SpeculativeState) SpeculativeState {
	return SpeculativeState{
		ID:            state.ID,
		RemainingWin:  state.RemainingWin,
		StartPC:       state.StartPC,
		Configuration: copyConfiguration(state.Configuration),
		CorrectPC:     state.CorrectPC,
	}
}

func copyExecutionPath(path ExecutionPath) ExecutionPath {
	newStack := make([]SpeculativeState, len(path.SpeculativeStack))
	for i, state := range path.SpeculativeStack {
		newStack[i] = copySpecState(state)
	}

	return ExecutionPath{
		CurrentConf:      copyConfiguration(path.CurrentConf),
		SpeculativeStack: newStack,
	}
}
