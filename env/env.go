package env

import (
	"os"
	"strconv"
)

func get(key string, defaultValue string) string {
	value, hasValue := os.LookupEnv(key)
	if hasValue {
		return value
	}
	return defaultValue
}

func TcpListenAddress() string {
	return get("TCP_ADDRESS", "localhost:9090")
}

func HttpListenAddress() string {
	return get("HTTP_ADDRESS", "localhost:9091")
}

func InternalMessageBufferSize() int {
	i, err := strconv.ParseInt(get("INTERNAL_MESSAGE_BUFFER_SIZE", "100"), 10, 32)
	if err != nil {
		return 100
	}

	return int(i)
}

func LogLevels() []string {
	level := get("LOG_LEVEL", "DEBUG")
	switch level {
	case "DEBUG":
		return []string{"ERROR", "WARN", "INFO", "DEBUG"}
	case "INFO":
		return []string{"ERROR", "WARN", "INFO"}
	case "WARN":
		return []string{"ERROR", "WARN"}
	case "ERROR":
		return []string{"ERROR"}
	}
	return []string{}
}
