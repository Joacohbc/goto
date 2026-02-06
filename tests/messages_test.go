package tests

import (
	"goto/src/core"
	"testing"
)

func TestNewMsg(t *testing.T) {
	msg := core.NewMsg(core.Info, "test message")
	if msg.Level != core.Info {
		t.Errorf("Expected level Info, got %v", msg.Level)
	}
	if msg.Content != "test message" {
		t.Errorf("Expected content 'test message', got '%s'", msg.Content)
	}
}

func TestNotifier(t *testing.T) {
	msgChan := make(chan core.Message, 1) // Buffer 1 to avoid blocking
	notifier := core.NewNotifier(msgChan)

	notifier.Notify(core.Success, "success msg")

	select {
	case msg := <-msgChan:
		if msg.Level != core.Success {
			t.Errorf("Expected level Success, got %v", msg.Level)
		}
		if msg.Content != "success msg" {
			t.Errorf("Expected content 'success msg', got '%s'", msg.Content)
		}
	default:
		t.Error("Did not receive message from channel")
	}

	// Test Info helper
	notifier.Info("info %s", "formatted")
	select {
	case msg := <-msgChan:
		if msg.Level != core.Info {
			t.Errorf("Expected level Info, got %v", msg.Level)
		}
		if msg.Content != "info formatted" {
			t.Errorf("Expected content 'info formatted', got '%s'", msg.Content)
		}
	default:
		t.Error("Did not receive message from channel")
	}

	// Test Success helper
	notifier.Success("success %s", "formatted")
	select {
	case msg := <-msgChan:
		if msg.Level != core.Success {
			t.Errorf("Expected level Success, got %v", msg.Level)
		}
		if msg.Content != "success formatted" {
			t.Errorf("Expected content 'success formatted', got '%s'", msg.Content)
		}
	default:
		t.Error("Did not receive message from channel")
	}

	// Test Warning helper
	notifier.Warning("warning %s", "formatted")
	select {
	case msg := <-msgChan:
		if msg.Level != core.Warning {
			t.Errorf("Expected level Warning, got %v", msg.Level)
		}
		if msg.Content != "warning formatted" {
			t.Errorf("Expected content 'warning formatted', got '%s'", msg.Content)
		}
	default:
		t.Error("Did not receive message from channel")
	}

	// Test Alert helper
	notifier.Alert("alert %s", "formatted")
	select {
	case msg := <-msgChan:
		if msg.Level != core.Alert {
			t.Errorf("Expected level Alert, got %v", msg.Level)
		}
		if msg.Content != "alert formatted" {
			t.Errorf("Expected content 'alert formatted', got '%s'", msg.Content)
		}
	default:
		t.Error("Did not receive message from channel")
	}

	// Test nil channel safety
	safeNotifier := core.NewNotifier(nil)
	safeNotifier.Info("should not panic")
}

func TestMessageHandler(t *testing.T) {
	var received []core.Message
	consume := func(msg core.Message) {
		received = append(received, msg)
	}

	mh := core.NewMessageHandler(consume)
	ch := mh.Channel()

	ch <- core.NewMsg(core.Info, "msg1")
	ch <- core.NewMsg(core.Warning, "msg2")

	mh.CloseAndWait()

	if len(received) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(received))
	}
	if received[0].Content != "msg1" {
		t.Errorf("Order mismatch or content wrong")
	}
}
