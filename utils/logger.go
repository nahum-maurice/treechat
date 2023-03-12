package utils

import (
	"fmt"
	"time"
)

type Logger struct {
	owner string
}

func NewLogger(owner string) *Logger {
	return &Logger{owner: owner}
}

func (l *Logger) Trace(s string) {
	formatted := fmt.Sprintf("[%v] %v TRACE ::: %v",
		l.owner,
		time.Now().Format("2006/01/02 15:04:05"),
		s)
	fmt.Println(formatted)
}

func (l *Logger) Debug(s string) {
	formatted := fmt.Sprintf("[%v] %v DEBUG ::: %v",
		l.owner,
		time.Now().Format("2006/01/02 15:04:05"),
		s)
	fmt.Println(formatted)
}

func (l *Logger) Info(s string) {
	formatted := fmt.Sprintf("[%v] %v INFO ::: %v",
		l.owner,
		time.Now().Format("2006/01/02 15:04:05"),
		s)
	fmt.Println(formatted)
}

func (l *Logger) Warn(s string) {
	formatted := fmt.Sprintf("[%v] %v WARNING ::: %v",
		l.owner,
		time.Now().Format("2006/01/02 15:04:05"),
		s)
	fmt.Println(formatted)
}

func (l *Logger) Error(s string) {
	formatted := fmt.Sprintf("[%v] %v ERROR ::: %v",
		l.owner,
		time.Now().Format("2006/01/02 15:04:05"),
		s)
	fmt.Println(formatted)
}

func (l *Logger) Fatal(s string) {
	formatted := fmt.Sprintf("[%v] %v FATAL ::: %v",
		l.owner,
		time.Now().Format("2006/01/02 15:04:05"),
		s)
	fmt.Println(formatted)
}
