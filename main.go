package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/loop_expander"
)

func main() {
	var inputFile string
	var outputFile string
	var unrollCount int

	flag.StringVar(&inputFile, "i", "", "入力アセンブリファイル")
	flag.StringVar(&outputFile, "o", "", "出力アセンブリファイル (指定しない場合は標準出力)")
	flag.IntVar(&unrollCount, "n", 2, "ループ展開回数")
	flag.Parse()

	if inputFile == "" {
		fmt.Println("入力ファイルを指定してください (-i オプション)")
		os.Exit(1)
	}

	if unrollCount <= 0 {
		fmt.Println("展開回数は正の整数である必要があります")
		os.Exit(1)
	}

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "入力ファイルのオープンに失敗しました: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	asm, err := assembler.ParseAsm(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "アセンブリコードのパースに失敗しました: %v\n", err)
		os.Exit(1)
	}

	expandedAsm, err := loop_expander.Loop_expander(asm, unrollCount)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ループ展開に失敗しました: %v\n", err)
		os.Exit(1)
	}

	// GenerateAsm を使用してアセンブリコードを文字列に変換
	output, err := assembler.GenerateAsm(expandedAsm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "アセンブリコードの生成に失敗しました: %v\n", err)
		os.Exit(1)
	}

	if outputFile != "" {
		err = os.WriteFile(outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "出力ファイルへの書き込みに失敗しました: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("ループ展開されたアセンブリコードを %s に書き込みました\n", outputFile)
	} else {
		fmt.Println(output)
	}
}
