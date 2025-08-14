package main

// LogMessages は、Mail IDとその関連するメッセージのリストを表します
type LogMessages struct {
	SortHint string           `json:"-"`
	MailId   string           `json:"mail_id"`
	Hostname string           `json:"hostname"`
	From     string           `json:"from"`
	To       string           `json:"to"`
	Status   string           `json:"status"`
	Messages []MessageContent `json:"messages"`
}

// MessageContent は、Postfixのログメッセージの内容を表します
type MessageContent struct {
	EDate  string            `json:"event_date"`
	ETime  string            `json:"event_time"`
	Params map[string]string `json:"paramaters"`
	RawMsg string            `json:"raw_message"`
}

// LogMessages にメッセージを追加します
func (lm *LogMessages) addMessage(msg MessageContent) {
	if lm.Messages != nil {
		lm.Messages = append(lm.Messages, msg)
	}
}
