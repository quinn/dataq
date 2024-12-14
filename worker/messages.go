package worker

import "fmt"

type Message struct {
	Type   string `json:"type"`
	Data   string `json:"data"`
	Closed bool   `json:"closed"`
	Done   bool   `json:"done"`
}

func sendError(messages chan<- Message, err error) {
	messages <- Message{
		Type: "error",
		Data: err.Error(),
	}
}

func sendErrorf(messages chan<- Message, format string, args ...interface{}) {
	sendError(messages, fmt.Errorf(format, args...))
}

func sendInfo(messages chan<- Message, info string) {
	messages <- Message{
		Type: "info",
		Data: info,
	}
}

func sendDone(messages chan<- Message) {
	messages <- Message{
		Done: true,
	}
}
