package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestConsoleLogger_Infof(t *testing.T) {
	logger := &ConsoleLogger{}

	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "simple string",
			format:   "hello world",
			args:     nil,
			expected: "hello world",
		},
		{
			name:     "formatted string",
			format:   "hello %s",
			args:     []interface{}{"world"},
			expected: "hello world",
		},
		{
			name:     "multiple arguments",
			format:   "%s: %d",
			args:     []interface{}{"count", 42},
			expected: "count: 42",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureStdout(func() {
				logger.Infof(tt.format, tt.args...)
			})

			if output != tt.expected {
				t.Errorf("got %q, want %q", output, tt.expected)
			}
		})
	}
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}
