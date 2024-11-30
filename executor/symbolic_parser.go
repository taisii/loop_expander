package executor

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// ParseSymbolicExpr parses a string expression into a SymbolicExpr.
func ParseSymbolicExpr(input string) (*SymbolicExpr, error) {
	tokens := tokenize(input)
	return parse(tokens)
}

// tokenize splits the input string into tokens.
func tokenize(input string) []string {
	replacements := []struct {
		old string
		new string
	}{
		{"(", " ( "}, {")", " ) "}, {"+", " + "}, {"-", " - "},
		{"*", " * "}, {"/", " / "}, {">", " > "}, {"<", " < "},
		{"=", " = "}, {"!=", " != "}, {">=", " >= "}, {"<=", " <= "},
	}
	for _, r := range replacements {
		input = strings.ReplaceAll(input, r.old, r.new)
	}
	return strings.Fields(input)
}

func parse(tokens []string) (*SymbolicExpr, error) {
	// ベースケース: トークンが1つだけの場合
	if len(tokens) == 1 {
		token := tokens[0]
		if isNumber(token) {
			// 数値トークンをリーフノードとして返す
			value, err := strconv.Atoi(token)
			if err != nil {
				return nil, fmt.Errorf("invalid number: %s", token)
			}
			return &SymbolicExpr{
				Op:       "value",
				Operands: []interface{}{value},
			}, nil
		} else if isIdentifier(token) {
			// 識別子トークンをリーフノードとして返す
			return &SymbolicExpr{
				Op:       "value",
				Operands: []interface{}{token},
			}, nil
		}
		return nil, fmt.Errorf("unexpected token: %s", token)
	}

	// 括弧で囲まれた式を処理
	if tokens[0] == "(" && tokens[len(tokens)-1] == ")" {
		// 括弧の内側を再帰的に解析
		return parse(tokens[1 : len(tokens)-1])
	}

	// 演算子が複数ある場合、最も優先順位の低い演算子を見つける
	index, err := findLowestPrecedenceOp(tokens)
	if err != nil {
		return nil, err
	}

	// 左右の部分式を再帰的に解析
	leftExpr, err := parse(tokens[:index])
	if err != nil {
		return nil, err
	}
	rightExpr, err := parse(tokens[index+1:])
	if err != nil {
		return nil, err
	}

	// 演算子と左右の部分式を結合して構文木を構築
	return &SymbolicExpr{
		Op:       tokens[index],
		Operands: []interface{}{*leftExpr, *rightExpr},
	}, nil
}

// Operator precedence map
var precedence = map[string]int{
	"==": 0, "!=": 0, // 等価、不等価比較（最低優先順位）
	"<": 0, "<=": 0, ">": 0, ">=": 0, // 大小比較

	"+": 1, "-": 1, // 加減算
	"*": 2, "/": 2, // 乗除算
	"^": 3, // 冪乗（右結合の演算子）
}

// findLowestPrecedenceOp finds the index of the operator with the lowest precedence
// in the given tokens, respecting left-to-right evaluation.
func findLowestPrecedenceOp(tokens []string) (int, error) {
	lowestPrec := int(^uint(0) >> 1) // 初期値は最大値（最小優先順位を探すため）
	lowestIdx := -1
	parenthesisLevel := 0 // 括弧のネストレベルを追跡

	for i, token := range tokens {
		switch token {
		case "(":
			// 括弧のネストを追跡
			parenthesisLevel++
		case ")":
			parenthesisLevel--
			if parenthesisLevel < 0 {
				return -1, fmt.Errorf("mismatched parentheses")
			}
		default:
			// 演算子で、かつ現在の括弧の外にいる場合のみ処理
			if prec, isOp := precedence[token]; isOp && parenthesisLevel == 0 {
				// 現在の演算子がより低い優先順位であれば更新
				if prec < lowestPrec {
					lowestPrec = prec
					lowestIdx = i
				}
			}
		}
	}

	if parenthesisLevel != 0 {
		return -1, fmt.Errorf("mismatched parentheses")
	}
	if lowestIdx == -1 {
		return -1, fmt.Errorf("no valid operator found")
	}
	return lowestIdx, nil
}

// isNumber checks if a token is a valid number (integer).
func isNumber(token string) bool {
	_, err := strconv.Atoi(token)
	return err == nil
}

// isIdentifier checks if a token is a valid identifier.
func isIdentifier(token string) bool {
	if len(token) == 0 {
		return false
	}
	for i, r := range token {
		if i == 0 {
			// 先頭はアルファベットまたはアンダースコア
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			// それ以降はアルファベット、数字、またはアンダースコア
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	return true
}
