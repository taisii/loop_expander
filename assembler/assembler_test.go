package assembler_test

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/taisii/go-project/assembler"
)

// TestLoadProgram: LoadProgramが正しく動作するかテスト
func TestLoadProgram(t *testing.T) {
	// テスト用の命令セットを定義
	ins := []assembler.Instruction{
		assembler.Label{Name: "start"},
		assembler.OpCode{Mnemonic: "mov", Operands: []string{"r1", "r2"}},
		assembler.Label{Name: "loop"},
		assembler.OpCode{Mnemonic: "jmp", Operands: []string{"start"}},
	}

	// アセンブラを初期化してプログラムをロード
	asm := assembler.NewAssembler()
	asm.LoadProgram(ins)

	// ラベルが正しいアドレスにマッピングされているか
	expectedLabels := map[string]int{"start": 0, "loop": 1}
	for label, addr := range expectedLabels {
		if asm.Labels[label] != addr {
			t.Errorf("Label %s expected to be at address %d, got %d", label, addr, asm.Labels[label])
		}
	}

	// プログラムが正しいアドレスと命令のペアになっているか

	expectedProgram := []assembler.Program{
		{Addr: 0, Inst: assembler.OpCode{Mnemonic: "mov", Operands: []string{"r1", "r2"}}},
		{Addr: 1, Inst: assembler.OpCode{Mnemonic: "jmp", Operands: []string{"0"}}},
	}
	for i, prog := range asm.Program {
		if prog.Addr != expectedProgram[i].Addr || prog.Inst.String() != expectedProgram[i].Inst.String() {
			t.Errorf("Expected program at index %d to be %+v, got %+v", i, expectedProgram[i], prog)
		}
	}
}

// TestShowProgram: ShowProgramの出力を確認
func TestShowProgram(t *testing.T) {
	// テスト用の命令セットを定義
	ins := []assembler.Instruction{
		assembler.Label{Name: "start"},
		assembler.OpCode{Mnemonic: "add", Operands: []string{"r1", "r2"}},
	}

	// アセンブラを初期化してプログラムをロード
	asm := assembler.NewAssembler()
	asm.LoadProgram(ins)

	// 出力をキャプチャ
	var buf bytes.Buffer
	originalStdout := os.Stdout // 現在のStdoutを保存
	r, w, _ := os.Pipe()        // Pipeを作成
	os.Stdout = w               // StdoutをPipeの書き込み側に差し替え

	// ShowProgramの出力を呼び出し
	go func() {
		asm.ShowProgram()
		w.Close() // 書き込み完了後に閉じる
	}()

	// Pipeの読み取り側から出力を取得
	io.Copy(&buf, r)
	os.Stdout = originalStdout // 元のStdoutを戻す

	// 期待される出力を定義
	expectedOutput := "  0: add r1, r2\nsym start = 0\n"
	if buf.String() != expectedOutput {
		t.Errorf("Expected output:\n%s\nGot:\n%s", expectedOutput, buf.String())
	}
}
