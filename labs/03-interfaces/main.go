package main

import "fmt"

// Logger is satisfied implicitly by any type with Log(string).
type Logger interface {
	Log(msg string)
}

type ConsoleLogger struct{}

func (ConsoleLogger) Log(msg string) { fmt.Println("[log]", msg) }

// Embedding promotes methods to the outer type.
type VerboseLogger struct {
	ConsoleLogger
	prefix string
}

func (v VerboseLogger) Log(msg string) {
	v.ConsoleLogger.Log(v.prefix + msg)
}

// Mock pattern for tests or demos.
type MockLogger struct {
	Messages []string
}

func (m *MockLogger) Log(msg string) {
	m.Messages = append(m.Messages, msg)
}

func notify(l Logger, text string) {
	l.Log("notify: " + text)
}

func main() {
	var l Logger = ConsoleLogger{}
	notify(l, "hello")

	v := VerboseLogger{prefix: ">> "}
	notify(v, "embedded type")

	mock := &MockLogger{}
	notify(mock, "recorded")
	fmt.Println("mock captured:", mock.Messages)
}