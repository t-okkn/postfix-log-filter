package main

import (
	"cmp"
	"compress/gzip"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// PATHがディレクトリかどうかを確認します
// true: ディレクトリ, false: ファイル
func isDirectory(path string) (bool, error) {
	info, err := os.Stat(path)

	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

// ファイルからログを読み込みます
// gzip圧縮されたファイルもサポートしています
func readFromFile(filePath string) ([]*LogMessages, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, fmt.Errorf("ファイルを開くことができません: %w", err)
	}

	defer file.Close()
	var data = []*LogMessages{}

	// gzip圧縮されたファイルかどうかを確認
	if strings.HasSuffix(filePath, ".gz") {
		gzr, err := gzip.NewReader(file)
		if err != nil {
			return data, fmt.Errorf("gzipファイルの展開に失敗: %w", err)
		}

		defer gzr.Close()

		if data, err = ParseLog(gzr); err != nil {
			return data, fmt.Errorf("gzipファイルの解析に失敗: %w", err)
		}

	} else {
		if data, err = ParseLog(file); err != nil {
			return data, fmt.Errorf("ファイルの解析に失敗: %w", err)
		}
	}

	return data, nil
}

// ディレクトリ内の全てのファイルからログを読み込みます
// ディレクトリ内のサブディレクトリは無視されます
func readFromDirectory(dirPath string) ([]*LogMessages, error) {
	files, err := os.ReadDir(dirPath)

	if err != nil {
		return nil, fmt.Errorf("ディレクトリを読み込むことができません: %w", err)
	}

	var allLogs []*LogMessages

	for _, file := range files {
		// ディレクトリは無視
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		logs, err := readFromFile(filePath)

		if err != nil {
			fmt.Printf("ファイル %s の読み込みに失敗しました: %v\n", filePath, err)
			continue
		}

		allLogs = append(allLogs, logs...)
	}

	slices.SortFunc(allLogs, func(a, b *LogMessages) int {
		return cmp.Compare(a.SortHint, b.SortHint)
	})

	return allLogs, nil
}
