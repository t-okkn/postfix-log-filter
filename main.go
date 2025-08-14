package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

var (
	PackageName string
	Version     string
	Revision    string
)

// main 関数は、コマンドライン引数を解析し、ログファイルの読み込みと結果の出力を行います
func main() {
	//変数でflagを定義します
	var (
		format  = flag.String("format", "json", "出力形式を指定します (json/csv) [default: json]")
		input   = flag.String("input", "", "入力元のファイル名・ディレクトリを指定します (パイプ以外は必須)")
		output  = flag.String("output", "", "出力先のファイル名を指定します (空の場合は標準出力) [default: 標準出力]")
		version = flag.Bool("version", false, "version情報を表示します")
	)

	flag.Parse()

	// versionフラグが指定された場合、バージョン情報を表示して終了
	if *version {
		fmt.Println(PackageName, "version:", Version, Revision)
		os.Exit(0)
	}

	var logs []*LogMessages

	// パイプで渡された内容かどうかを確認
	if term.IsTerminal(int(os.Stdin.Fd())) {
		if *input == "" {
			fmt.Println("入力元のファイル名・ディレクトリを指定してください")
			fmt.Println()
			flag.Usage()
			os.Exit(127)
		}

		isDir, err := isDirectory(*input)

		if err != nil {
			fmt.Printf("入力されたPATHが不正です: %v\n", err)
			os.Exit(126)
		}

		if isDir {
			logs, err = readFromDirectory(*input)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		} else {
			logs, err = readFromFile(*input)

			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

	} else {
		// 【パイプ】標準入力からログを解析
		var err error
		logs, err = ParseLog(os.Stdin)

		if err != nil {
			fmt.Printf("ログの解析に失敗しました: %v\n", err)
			os.Exit(1)
		}

		if len(logs) == 0 {
			fmt.Println("正常終了：メール送受信ログが見つかりませんでした")
			os.Exit(0)
		}
	}

	// 出力形式を決定
	var outputWriter io.Writer
	if *output == "" {
		outputWriter = os.Stdout

	} else {
		file, err := os.Create(*output)

		if err != nil {
			fmt.Printf("出力ファイルの作成に失敗しました: %v\n", err)
			os.Exit(125)
		}

		defer file.Close()
		outputWriter = file
	}

	// 出力形式に応じてエクスポート
	if strings.ToLower(*format) == "csv" {
		if err := exportCSV(logs, outputWriter); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}

	} else {
		// デフォルトはJSON形式
		if err := exportJSON(logs, outputWriter); err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}
}
