package main

import (
	"fmt"
	"os"

	"github.com/taisii/go-project/engine"
)

var dirPath string = "tests/"

func main() {
	// 読み込むディレクトリのパスを指定

	// ディレクトリ内のファイル一覧を取得
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// 各ファイルの内容を読み込む
	for _, file := range files {
		if file.IsDir() {
			// サブディレクトリはスキップ
			continue
		}

		filePath := dirPath + file.Name()
		content, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
			continue
		}

		// ファイルの内容を出力
		fmt.Printf("Content of %s:\n", file.Name())
		engine.Execute(content,10)
		fmt.Println("------------")
	}
}
