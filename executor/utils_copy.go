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
