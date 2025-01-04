package assembler

func CopyInstructionValue(inst *Instruction) *Instruction {
	newInst := *inst
	return &newInst
}

func CopyOpCodeValue(op *OpCode) *OpCode {
	newOp := *op
	return &newOp
}

func CopyLabelValue(label *Label) *Label {
	newLabel := *label
	return &newLabel
}

func CopyAssembler(asm *Assembler) *Assembler {
	newAsm := &Assembler{
		Program: make([]Instruction, len(asm.Program)),
		Labels:  make(map[string]int, len(asm.Labels)),
	}

	// Instructionのコピー
	for i, inst := range asm.Program {
		newAsm.Program[i] = *CopyInstructionValue(&inst)
	}

	// Labelsマップのコピー
	for k, v := range asm.Labels {
		newAsm.Labels[k] = v
	}

	return newAsm
}
