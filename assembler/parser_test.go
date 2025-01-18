package assembler_test

import (
	"os"
	"testing"

	"github.com/taisii/go-project/assembler"
)

func TestParseAsm(t *testing.T) {
	// テストケース
	testCases := []struct {
		filename         string // テストファイル名
		expectedAssembly assembler.Assembler
	}{
		{
			filename: "../tests/test1.muasm", // ファイル名は適宜変更してください
			expectedAssembly: assembler.Assembler{
				Labels: map[string]int{
					"Loop": 2,
				},
				Program: []assembler.Instruction{
					{0, assembler.OpCode{"<-", []string{"x", "5"}}},
					{1, assembler.OpCode{"<-", []string{"w", "0"}}},
					{2, assembler.OpCode{"<-", []string{"w", "w+x"}}},
					{3, assembler.OpCode{"<-", []string{"x", "x-1"}}},
					{4, assembler.OpCode{"<-", []string{"y", "x=0"}}},
					{5, assembler.OpCode{"beqz", []string{"y", "Loop"}}},
				},
			},
		},
		{
			filename: "../tests/test2.muasm", // ファイル名は適宜変更してください
			expectedAssembly: assembler.Assembler{
				Labels: map[string]int{
					"L6":  6,
					"L7":  7,
					"L10": 10,
					"End": 11,
				},
				Program: []assembler.Instruction{
					{0, assembler.OpCode{"load", []string{"x", "0"}}},
					{1, assembler.OpCode{"load", []string{"v", "1"}}},
					{2, assembler.OpCode{"load", []string{"w", "2"}}},
					{3, assembler.OpCode{"beqz", []string{"x", "L6"}}},
					{4, assembler.OpCode{"load", []string{"y", "v"}}},
					{5, assembler.OpCode{"jmp", []string{"L7"}}},
					{6, assembler.OpCode{"store", []string{"y", "w"}}},
					{7, assembler.OpCode{"beqz", []string{"x", "L10"}}},
					{8, assembler.OpCode{"store", []string{"y", "w"}}},
					{9, assembler.OpCode{"jmp", []string{"End"}}},
					{10, assembler.OpCode{"load", []string{"y", "v"}}},
				},
			},
		},
		{
			filename: "../tests/test3.muasm",
			expectedAssembly: assembler.Assembler{
				Labels: map[string]int{
					"L3":  3,
					"L10": 5,
				},
				Program: []assembler.Instruction{
					{0, assembler.OpCode{"<-", []string{"x", "in>=bound"}}},
					{1, assembler.OpCode{"beqz", []string{"x", "L3"}}},
					{2, assembler.OpCode{"jmp", []string{"L10"}}},
					{3, assembler.OpCode{"load", []string{"secret", "in"}}},
					{4, assembler.OpCode{"load", []string{"z", "secret"}}},
				},
			},
		},
				{
			filename: "../tests/test4.muasm",
			expectedAssembly: assembler.Assembler{
				Labels: map[string]int{
					"End":  5,
				},
				Program: []assembler.Instruction{
					{0, assembler.OpCode{"<-", []string{"x", "v<y"}}},
					{1, assembler.OpCode{"beqz", []string{"x", "End"}}},
					{2, assembler.OpCode{"spbarr", []string{""}}},
					{3, assembler.OpCode{"load", []string{"v", "v"}}},
					{4, assembler.OpCode{"load", []string{"v", "v"}}},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			file, err := os.Open(tc.filename)
			if err != nil {
				t.Fatalf("ファイルを開けませんでした: %v", err)
			}
			defer file.Close()

			got, err := assembler.ParseAsm(file)
			if err != nil {
				t.Fatalf("parseAsmエラー: %v", err)
			}

			// ラベルのテスト
			if len(got.Labels) != len(tc.expectedAssembly.Labels) {
				t.Errorf("ラベルの数が異なります。got: %d, want: %d", len(got.Labels), len(tc.expectedAssembly.Labels))
			}
			for label, addr := range tc.expectedAssembly.Labels {
				gotAddr, ok := got.Labels[label]
				if !ok {
					t.Errorf("ラベル %s が見つかりません", label)
				} else if gotAddr != addr {
					t.Errorf("ラベル %s のアドレスが異なります。got: %d, want: %d", label, gotAddr, addr)
				}
			}

			// プログラムのテスト
			if len(got.Program) != len(tc.expectedAssembly.Program) {
				t.Errorf("プログラムの長さが異なります。got: %d, want: %d", len(got.Program), len(tc.expectedAssembly.Program))
			}
			for i, wantInst := range tc.expectedAssembly.Program {
				gotInst := got.Program[i]
				if gotInst.Addr != wantInst.Addr {
					t.Errorf("アドレスが異なります。index: %d, got: %d, want: %d", i, gotInst.Addr, wantInst.Addr)
				}
				if gotInst.OpCode.String() != wantInst.OpCode.String() {
					t.Errorf("命令が異なります。index: %d, got: %s, want: %s", i, gotInst.OpCode.String(), wantInst.OpCode.String())
				}
			}
		})
	}
}
