package executor

// Configuration structure
type Configuration struct {
	PC        int                    // Program Counter
	Registers map[string]interface{} // General-purpose registers (can hold symbolic or concrete values)
	Memory    map[int]interface{}    // Memory (address to value, symbolic or concrete)
	Trace     Trace
	StepCount int
}

// SymbolicExpr represents a symbolic expression.
type SymbolicExpr struct {
	Op       string        // Operator ("+", "-", ">", etc.)
	Operands []interface{} // Operands (can be integers, strings, or nested SymbolicExpr)
}

type Trace struct {
	Observations []Observation // 実行過程の観測データのリスト
	PathCond     SymbolicExpr  // このトレースのパス条件（シンボリック形式）
}

type Observation struct {
	PC        int               // プログラムカウンタ（実行中の命令の位置）
	Type      ObsType           // 観測タイプ: load, store, pc, start, rollback, commit
	Address   interface{}       // メモリアクセスの場合のアドレス（シンボリック形式）
	Value     interface{}       // 値の読み取りや書き込みの内容（シンボリック形式）
	SpecState *SpeculativeState // スペキュレーション状態（該当する場合）
}

type ObsType string

const (
	ObsTypeLoad     ObsType = "load"     // メモリ読み取り
	ObsTypeStore    ObsType = "store"    // メモリ書き込み
	ObsTypePC       ObsType = "pc"       // プログラムカウンタの変更
	ObsTypeStart    ObsType = "start"    // スペキュレーション開始
	ObsTypeRollback ObsType = "rollback" // スペキュレーション取り消し
	ObsTypeCommit   ObsType = "commit"   // スペキュレーションのコミット
)

// SpeculativeState represents an individual speculative execution state.
type SpeculativeState struct {
	ID           int           // Unique identifier for the speculative state
	RemainingWin int           // Remaining speculative window size (N)
	StartPC      int           // Program Counter (PC) where speculation started
	InitialConf  Configuration // Configuration when speculation started
	CorrectPC    int           // Correct branch PC to commit if speculation is valid
}

// ExecutionState represents the overall execution state during speculative execution.
type ExecutionState struct {
	Counter     int                // Counter for generating unique IDs for speculative states
	CurrentConf Configuration      // Current non-speculative configuration
	Speculative []SpeculativeState // Stack of speculative states
}
