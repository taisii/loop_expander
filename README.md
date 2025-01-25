# SPECTECTORの単純ループへの拡張

本リポジトリは、卒業論文「SPECTECTORの単純ループへの拡張」で提案されたアルゴリズムの実装です。

**注意:** この論文は学術雑誌や会議等で公開されたものではなく、卒業論文として提出されたものです。

## 概要

このアルゴリズムは[SPECTECTOR](https://github.com/spectector/spectector)の前処理として実装されたものです。SPECTECTORのアルゴリズムがループを含むプログラムに対して実行が停止しないことがある問題に対する解決策として作られています。

## 実装

この実装は、Go言語(version 1.23.3)で記述されており、Goの環境構築がされていることが前提となっています。

## 再現実験

論文で報告されている実験結果を再現するための手順は[loop_expander_testリポジトリ](https://github.com/taisii/loop_expander_test) を参照してください。

## 実行環境

*   Go バージョン: 1.23.3
*   OS: macOS Sequoia 15.1.1