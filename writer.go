package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// ログメッセージのスライスをCSV形式で出力します
func exportCSV(logs []*LogMessages, output io.Writer) error {
	w := csv.NewWriter(output)

	// ヘッダー行を追加
	w.Write([]string{
		"Hostname",
		"MailID",
		"Sequence",
		"EventDate",
		"EventTime",
		"From",
		"To",
		"Status",
		"RawMessage",
	})

	for _, log := range logs {
		counter := 1

		for _, msg := range log.Messages {
			w.Write([]string{
				log.Hostname,
				log.MailId,
				fmt.Sprintf("%d", counter),
				msg.EDate,
				msg.ETime,
				log.From,
				log.To,
				log.Status,
				msg.RawMsg,
			})

			counter += 1
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return fmt.Errorf("CSVの出力に失敗: %w", err)
	}

	return nil
}

// ログメッセージのスライスをJSON形式で出力します
func exportJSON(logs []*LogMessages, output io.Writer) error {
	encoder := json.NewEncoder(output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", strings.Repeat(" ", 2))

	if err := encoder.Encode(logs); err != nil {
		return fmt.Errorf("JSONの出力に失敗: %w", err)
	}

	return nil
}