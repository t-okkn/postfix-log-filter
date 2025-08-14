package main

import (
	"bufio"
	"io"
	"strings"
	"time"
)

// Postfixのログを解析します
func ParseLog(logdata io.Reader) ([]*LogMessages, error) {
	logs := make(map[string]*LogMessages)
	hostname := ""
	scanner := bufio.NewScanner(logdata)

	for scanner.Scan() {
		line := scanner.Text()

		// 空行はスキップ
		if len(line) == 0 {
			continue
		}

		// ホスト名を取得
		if hostname == "" {
			hostname = parseHostname(line)
		}

		// 行を解析して、Mail IDとMessageContentを取得
		mailId, msgContent := parseLine(line)

		// Mail IDが空の場合はスキップ
		if mailId == "" {
			continue
		}

		// ログメッセージを取得または作成
		logMessage, exists := logs[mailId]

		// 存在しない場合は新しいLogMessagesを作成
		if !exists {
			logMessage = &LogMessages{
				SortHint: "",
				MailId:   mailId,
				Hostname: hostname,
				From:     "",
				To:       "",
				Status:   "",
				Messages: []MessageContent{},
			}

			logs[mailId] = logMessage
		}

		// メッセージを追加
		logMessage.addMessage(msgContent)

		if logMessage.SortHint == "" {
			// SortHintを設定（例: "1001123456123456789AB" の形式）
			d := strings.Trim(msgContent.EDate, "-")
			t := strings.Trim(msgContent.ETime, ":")
			logMessage.SortHint = d + t + mailId
		}

		// From, To, Status の情報を必要に応じて logMessage に設定
		if logMessage.From == "" {
			from, fromExsists := msgContent.Params["from"]

			if fromExsists {
				logMessage.From = from
			}
		}

		if logMessage.To == "" {
			to, toExists := msgContent.Params["to"]
			origTo, origToExists := msgContent.Params["orig_to"]

			if origToExists {
				logMessage.To = origTo

			} else if toExists {
				logMessage.To = to
			}
		}

		if logMessage.Status == "" {
			status, statusExists := msgContent.Params["status"]

			if statusExists {
				logMessage.Status = status
			}
		}
	}

	return toSlice(logs), nil
}

// 行を解析します
func parseLine(line string) (string, MessageContent) {
	// 行の長さが短い場合は、適切な情報が不足しているため、空の値を返します
	if len(line) <= 16 {
		return "", MessageContent{}
	}

	elements := strings.Fields(line[16:])

	// 要素が3つ未満の場合は、適切な情報が不足しているため、空の値を返します
	if len(elements) < 3 {
		return "", MessageContent{}
	}

	// Mail IDの位置にMail IDが含まれていなければ空の値を返します
	if !strings.HasSuffix(elements[2], ":") {
		return "", MessageContent{}
	}

	// 統計情報の行は無視
	if elements[2] == "statistics:" {
		return "", MessageContent{}
	}

	mailId := strings.TrimSuffix(elements[2], ":")
	dt, _ := time.Parse("Jan 2 15:04:05", line[:15])
	mc := MessageContent{
		EDate:  dt.Format("01-02"),
		ETime:  dt.Format("15:04:05"),
		Params: make(map[string]string),
		RawMsg: strings.Join(elements[1:], " "),
	}

	if elements[2] == "warning:" && strings.HasSuffix(elements[3], ":") {
		// 警告メッセージの場合、Mail IDは3番目の要素に含まれていることがある
		mailId = strings.TrimSuffix(elements[3], ":")
		mc.Params["warning"] = strings.Join(elements[4:], " ")
	}

	commentStart := -1
	relationship := -1

	// パラメータを解析
	for i := 3; i < len(elements); i++ {
		if strings.Contains(elements[i], "=") {
			parts := strings.Split(elements[i], "=")

			if len(parts) == 2 {
				key := parts[0]
				value := strings.Trim(parts[1], "<>,")

				mc.Params[key] = value
			}

		} else if strings.HasPrefix(elements[i], "(") {
			commentStart = i
			break

		} else if strings.Contains(elements[i], "notification:") {
			relationship = i + 1
			break
		}
	}

	if commentStart != -1 {
		comment := strings.Join(elements[commentStart:], " ")
		mc.Params["comment"] = strings.TrimLeft(comment, "(")[:len(comment)-2]
	}

	if relationship != -1 {
		mc.Params["relationship"] = elements[relationship]
	}

	return mailId, mc
}

// ホスト名を抽出します
func parseHostname(line string) string {
	if len(line) <= 16 {
		return ""
	}

	return strings.Fields(line[16:])[0]
}

// LogMessages のマップをスライスに変換します
func toSlice(logs map[string]*LogMessages) []*LogMessages {
	slice := make([]*LogMessages, 0, len(logs))

	for _, log := range logs {
		if log != nil {
			slice = append(slice, log)
		}
	}

	return slice
}
