package core

import (
	"fmt"
	"sync"
)

// MsgLevel defines the severity or type of the message.
type MsgLevel int

const (
	Info MsgLevel = iota
	Success
	Warning
	Alert
)

// Message represents a status update to be sent to the presentation layer.
type Message struct {
	Level   MsgLevel
	Content string
}

// NewMsg creates a new Message instance.
func NewMsg(level MsgLevel, content string) Message {
	return Message{
		Level:   level,
		Content: content,
	}
}

// Notifier encapsulates a message channel to provide helper methods for sending updates.
type Notifier struct {
	msgChan chan<- Message
}

// NewNotifier creates a new Notifier instance.
func NewNotifier(msgChan chan<- Message) *Notifier {
	return &Notifier{msgChan: msgChan}
}

// Notify sends a raw message if the channel is not nil.
func (n *Notifier) Notify(level MsgLevel, content string) {
	if n.msgChan != nil {
		n.msgChan <- NewMsg(level, content)
	}
}

// Info sends an Info message with formatting.
func (n *Notifier) Info(format string, args ...interface{}) {
	n.Notify(Info, fmt.Sprintf(format, args...))
}

// Success sends a Success message with formatting.
func (n *Notifier) Success(format string, args ...interface{}) {
	n.Notify(Success, fmt.Sprintf(format, args...))
}

// Warning sends a Warning message with formatting.
func (n *Notifier) Warning(format string, args ...interface{}) {
	n.Notify(Warning, fmt.Sprintf(format, args...))
}

// Alert sends an Alert message with formatting.
func (n *Notifier) Alert(format string, args ...interface{}) {
	n.Notify(Alert, fmt.Sprintf(format, args...))
}

// MessageHandler encapsulates the logic for consuming messages asynchronously.
type MessageHandler struct {
	msgChan chan Message
	wg      sync.WaitGroup
}

// NewMessageHandler creates and starts a new MessageHandler.
// The consumeFunc is called for each message received.
func NewMessageHandler(consumeFunc func(Message)) *MessageHandler {
	mh := &MessageHandler{
		msgChan: make(chan Message),
	}
	mh.wg.Add(1)
	go func() {
		defer mh.wg.Done()
		for msg := range mh.msgChan {
			consumeFunc(msg)
		}
	}()
	return mh
}

// Channel returns the channel to send messages to.
func (mh *MessageHandler) Channel() chan<- Message {
	return mh.msgChan
}

// CloseAndWait closes the channel and waits for processing to complete.
func (mh *MessageHandler) CloseAndWait() {
	close(mh.msgChan)
	mh.wg.Wait()
}
