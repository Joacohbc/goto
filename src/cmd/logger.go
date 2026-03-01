package cmd

import "fmt"

// ConsoleLogger implementa core.Logger para imprimir en consola standard
type ConsoleLogger struct{}

func (l *ConsoleLogger) Infof(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
