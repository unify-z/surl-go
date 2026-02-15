package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	DEBUG   = 0
	INFO    = 1
	WARNING = 2
	ERROR   = 3
	FATAL   = 4

	logLevel int

	consoleLogger *log.Logger
	fileLogger    *log.Logger
	fileHandle    *os.File
)

const (
	ansiReset  = "\033[0m"
	ansiRed    = "\033[31m"
	ansiYellow = "\033[33m"
	ansiGreen  = "\033[32m"
	ansiGray   = "\033[90m"
)

func GetLogLevelInt(levelStr string) int {
	switch levelStr {
	case "info":
		return INFO
	case "warning":
		return WARNING
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO

	}

}

func InitLogger(level int, filePath string) error {
	logLevel = level
	consoleLogger = log.New(os.Stdout, "", log.LstdFlags)

	if fileHandle != nil {
		_ = fileHandle.Close()
		fileHandle = nil
		fileLogger = nil
	}

	if filePath != "" {
		f, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		fileHandle = f
		fileLogger = log.New(fileHandle, "", log.LstdFlags)
	}

	return nil
}

func Writer() io.Writer {
	if fileHandle != nil {
		return fileHandle
	}
	return os.Stdout
}

func write(tag, color, format string, v ...any) {
	msg := fmt.Sprintf(format, v...)

	if consoleLogger != nil {
		consoleLogger.Printf("%s[%s] %s%s", color, tag, msg, ansiReset)
	}

	if fileLogger != nil {
		fileLogger.Printf("[%s] %s", tag, msg)
	}
}

func Debug(format string, v ...any) {
	if logLevel <= DEBUG {
		write("DEBUG", ansiGray, format, v...)
	}
}

func Info(format string, v ...any) {
	if logLevel <= INFO {
		write("INFO", ansiGreen, format, v...)
	}
}

func Warning(format string, v ...any) {
	if logLevel <= WARNING {
		write("WARNING", ansiYellow, format, v...)
	}
}

func Error(format string, v ...any) {
	if logLevel <= ERROR {
		write("ERROR", ansiRed, format, v...)
	}
}

func Fatal(format string, v ...any) {
	if fileLogger != nil {
		fileLogger.Printf("[FATAL] %s", fmt.Sprintf(format, v...))
	}

	if consoleLogger != nil {
		consoleLogger.Fatalf("%s[FATAL] %s%s", ansiRed, fmt.Sprintf(format, v...), ansiReset)
	}

	os.Exit(1)
}

func Close() error {
	if fileHandle != nil {
		err := fileHandle.Close()
		fileHandle = nil
		fileLogger = nil
		return err
	}
	return nil
}
