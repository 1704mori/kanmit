package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	successLogger *log.Logger
}

var (
	Colorize        = color.New(color.FgYellow).SprintfFunc()
	ColorizeSuccess = color.New(color.FgGreen).SprintfFunc()
	ColorizeError   = color.New(color.FgRed).SprintfFunc()
)

func NewLogger() *Logger {
	return &Logger{
		infoLogger:    log.New(os.Stdout, Colorize("INFO: "), 0),
		errorLogger:   log.New(os.Stderr, ColorizeError("ERROR: "), 0),
		successLogger: log.New(os.Stdout, ColorizeSuccess("SUCCESS: "), 0),
	}
}

func (l *Logger) Info(message string) {
	l.infoLogger.Println(message)
}

func (l *Logger) Success(message string) {
	l.successLogger.Println(message)
}

func (l *Logger) Error(message string) {
	l.errorLogger.Println(message)
}

func (l *Logger) AnimateLoading(message string, done chan struct{}) {
	loadingSymbols := `-\|/`
	i := 0

	for {
		select {
		case <-done:
			return
		default:
			fmt.Printf("\r%s %s", message, string(loadingSymbols[i]))
			time.Sleep(100 * time.Millisecond)
			i = (i + 1) % len(loadingSymbols)
		}
	}
}

func (l *Logger) ClearLoading() {
	fmt.Print("\r\x1b[K")
}
