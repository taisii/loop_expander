package assembler_test

import (
	"testing"

	"github.com/taisii/go-project/assembler"
)

func TestGenerateAsm(t *testing.T) {
	tests := []struct {
		name    string
		input   *assembler.Assembler
		want    string
		wantErr bool
	}{
		{
			name: "基本的な命令とラベル",
			input: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "SET", Operands: []string{"R1", "10"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "JMP", Operands: []string{"loop"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "HALT", Operands: []string{}}},
				},
				Labels: map[string]int{
					"start": 0,
					"loop":  2,
				},
			},
			want: `start:
SET R1, 10
JMP loop
loop:
HALT
`,
			wantErr: false,
		},
		{
			name: "代入命令",
			input: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "<-", Operands: []string{"x", "R1"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "<-", Operands: []string{"R2", "x"}}},
				},
				Labels: map[string]int{},
			},
			want: `x <- R1
R2 <- x
`,
			wantErr: false,
		},
		{
			name: "オペランドなしの命令",
			input: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "NOP", Operands: []string{}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "RET", Operands: []string{}}},
				},
				Labels: map[string]int{},
			},
			want: `NOP
RET
`,
			wantErr: false,
		},
		{
			name: "複数のラベル",
			input: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "SET", Operands: []string{"R1", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "JMP", Operands: []string{"middle"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "HALT", Operands: []string{}}},
				},
				Labels: map[string]int{
					"start":  0,
					"middle": 2,
				},
			},
			want: `start:
SET R1, 0
JMP middle
middle:
HALT
`,
			wantErr: false,
		},
		{
			name: "Loop expander",
			input: &assembler.Assembler{
				Program: []assembler.Instruction{
					{Addr: 0, OpCode: assembler.OpCode{Mnemonic: "load", Operands: []string{"x", "0"}}},
					{Addr: 1, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 2, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "LoopStart_0"}}},
					{Addr: 3, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 4, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 5, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "LoopStart_1"}}},
					{Addr: 6, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
					{Addr: 7, OpCode: assembler.OpCode{Mnemonic: "add", Operands: []string{"x", "1"}}},
					{Addr: 8, OpCode: assembler.OpCode{Mnemonic: "beqz", Operands: []string{"x", "LoopStart_2"}}},
					{Addr: 9, OpCode: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"programEnd"}}},
				},
				Labels: map[string]int{
					"LoopStart":   1,
					"LoopStart_0": 4,
					"LoopStart_1": 7,
					"LoopStart_2": 10,
					"programEnd":  10,
				},
			},
			want: `load x, 0
LoopStart:
add x, 1
beqz x, LoopStart_0
jmp programEnd
LoopStart_0:
add x, 1
beqz x, LoopStart_1
jmp programEnd
LoopStart_1:
add x, 1
beqz x, LoopStart_2
jmp programEnd
LoopStart_2:
programEnd:
`,
		},
		{
			name: "spbarr",
			input: &assembler.Assembler{
				Labels: map[string]int{
					"End": 5,
				},
				Program: []assembler.Instruction{
					{0, assembler.OpCode{"<-", []string{"x", "v<y"}}},
					{1, assembler.OpCode{"beqz", []string{"x", "End"}}},
					{2, assembler.OpCode{"spbarr", []string{""}}},
					{3, assembler.OpCode{"load", []string{"v", "v"}}},
					{4, assembler.OpCode{"load", []string{"v", "v"}}},
				},
			},
			want: `x <- v<y
beqz x, End
spbarr
load v, v
load v, v
End:
`,
			wantErr: false,
		},
		{
			name: "空のAssembler",
			input: &assembler.Assembler{
				Program: []assembler.Instruction{},
				Labels:  map[string]int{},
			},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := assembler.GenerateAsm(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAsm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateAsm() got = \n%v\n, want \n%v", got, tt.want)
			}
		})
	}
}
