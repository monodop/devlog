package log

import "fmt"

func log(level string, message string) {
	fmt.Printf("[%s] %s\n", level, message)
}

func Debug(message string, args ...interface{}) {
	log("DEBUG", fmt.Sprintf(message, args...))
}

func Info(message string, args ...interface{}) {
	log("INFO", fmt.Sprintf(message, args...))
}

func Warn(message string, args ...interface{}) {
	log("WARN", fmt.Sprintf(message, args...))
}

func Error(message string, args ...interface{}) {
	log("ERROR", fmt.Sprintf(message, args...))
}

func Exception(err error) {
	Error("%s", err)
}
