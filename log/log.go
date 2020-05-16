package log

import (
	"fmt"
	"strings"

	"github.com/monodop/devlog/env"
)

func log(level string, message string) {
	allowedLevels := env.LogLevels()
	allowed := false
	for _, allowedLevel := range allowedLevels {
		if strings.ToUpper(allowedLevel) == strings.ToUpper(level) {
			allowed = true
			break
		}
	}
	if allowed {
		fmt.Printf("[%s] %s\n", level, message)
	}
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
