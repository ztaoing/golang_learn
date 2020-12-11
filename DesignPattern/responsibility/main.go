package main

import "fmt"

//责任链模式：好处是解耦，如果不用责任链模式把处理分层，
//那么就需要一个硕大的函数来处理各种判断，那么这个函数就会很臃肿里面包含各种if...else...
const (
	Debug = iota
	Info
	Warning
	Error
	Fatal
)

type Logger interface {
	PrintLog(level int, msg string)
}

var ()

type BaseLogger struct {
	next Logger
}

func (b *BaseLogger) SetNext(logger Logger) {
	b.next = logger
}

type DebugLogger struct {
	BaseLogger
}

func (d *DebugLogger) PrintLog(level int, msg string) {
	if level == Debug {
		fmt.Printf("[DEBUG] %s\n", msg)
	} else {
		fmt.Printf("ignore [DEBUG]\n")

	}
}

type InfoLogger struct {
	BaseLogger
}

func (i *InfoLogger) PrintLog(level int, msg string) {
	if level == Info {
		fmt.Printf("[INFO] %s\n", msg)
	} else {
		fmt.Printf("ignore [INFO]\n")

	}
}

type ErrorLogger struct {
	BaseLogger
}

func (e *ErrorLogger) PrintLog(level int, msg string) {
	if level == Error {
		fmt.Printf("[ERROR %s\n]", msg)
	} else {
		fmt.Printf("ignore [ERROR]\n")

	}
}

func main() {
	errorLogger := &ErrorLogger{}
	infoLogger := &InfoLogger{}
	debugLogger := &DebugLogger{}

	infoLogger.SetNext(infoLogger)
	errorLogger.SetNext(errorLogger)
	debugLogger.SetNext(debugLogger)

	errorLogger.PrintLog(Info, "info")
	errorLogger.PrintLog(Debug, "debug")
}
