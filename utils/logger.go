package utils

import (
	"fmt"
)

type Logger struct {
	owner     string
	formatter *Formatter
}

func NewLogger(owner string) *Logger {
	return &Logger{
		owner: owner, 
		formatter: NewFormatter("System"),
	}
}

func (l *Logger) Trace(s string) {
	formatted := l.formatter.Log(s, "TRACE")
	fmt.Println(formatted)
}

func (l *Logger) Debug(s string) {
	formatted := l.formatter.Log(s, "DEBUG")
	fmt.Println(formatted)
}

func (l *Logger) Info(s string) {
	formatted := l.formatter.Log(s, "INFO")
	fmt.Println(formatted)
}

func (l *Logger) Warn(s string) {
	formatted := l.formatter.Log(s, "WARNING")
	fmt.Println(formatted)
}

func (l *Logger) Error(s string) {
	formatted := l.formatter.Log(s, "ERROR")
	fmt.Println(formatted)
}

func (l *Logger) Fatal(s string) {
	formatted := l.formatter.Log(s, "FATAL")
	fmt.Println(formatted)
}
