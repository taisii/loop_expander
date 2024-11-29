package executor_test

import (
	"testing"

	"github.com/taisii/go-project/assembler"
	"github.com/taisii/go-project/executor"
)

type TestCase struct {
	Name            string                   // テストケースの名前
	Program         []assembler.OpCode       // プログラム
	InitialConfig   *executor.Configuration  // 初期状態
	MaxSteps        int                      // 最大ステップ数
	ExpectedConfigs []executor.Configuration // 期待される終了状態
	ExpectError     bool                     // エラーを期待するか
}

func RunTestCase(t *testing.T, testCases []TestCase, executeFunc func([]assembler.OpCode, *executor.Configuration, int) ([]*executor.Configuration, error)) {
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			// 実行関数を呼び出し
			finalConfigs, err := executeFunc(testCase.Program, testCase.InitialConfig, testCase.MaxSteps)

			// エラーの確認
			if testCase.ExpectError {
				if err == nil {
					t.Errorf("Test case '%s' failed: expected an error but got none", testCase.Name)
				}
			} else {
				if err != nil {
					t.Errorf("Test case '%s' failed: unexpected error: %v", testCase.Name, err)
					return
				}

				// 終了状態の数が期待通りか確認
				if len(finalConfigs) != len(testCase.ExpectedConfigs) {
					t.Errorf("Test case '%s' failed: expected %d final configurations, but got %d",
						testCase.Name, len(testCase.ExpectedConfigs), len(finalConfigs))
					return
				}

				// 各終了状態の比較
				for i, expectedConfig := range testCase.ExpectedConfigs {
					if !executor.CompareConfiguration(expectedConfig, *finalConfigs[i]) {
						difference := executor.FormatConfigDifferences(expectedConfig, *finalConfigs[i])
						t.Errorf("Test case '%s' failed: configuration %d did not match expected configuration.\n%s",
							testCase.Name, i+1, difference)
						executor.PrintConfiguration(*finalConfigs[i])
					}
				}
			}
		})
	}
}
