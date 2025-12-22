package logger

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	mu   sync.Mutex
	file *os.File
	log  *log.Logger
}

type entry struct {
	Timestamp string                 `json:"ts"`
	Level     string                 `json:"level"`
	Message   string                 `json:"msg"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

var global = &Logger{log: log.New(os.Stdout, "", 0)}

func Init(path string) error {
	if path == "" {
		path = filepath.Join("logs", "app.jsonl")
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	global.mu.Lock()
	if global.file != nil {
		_ = global.file.Close()
	}
	global.file = f
	global.log = log.New(f, "", 0)
	global.mu.Unlock()
	return nil
}

func Close() error {
	global.mu.Lock()
	defer global.mu.Unlock()
	if global.file == nil {
		return errors.New("logger not initialized")
	}
	err := global.file.Close()
	global.file = nil
	return err
}

func logEntry(level, msg string, fields map[string]interface{}) {
	global.mu.Lock()
	defer global.mu.Unlock()
	payload := entry{
		Timestamp: time.Now().Format(time.RFC3339Nano),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}
	data, err := json.Marshal(payload)
	if err != nil {
		global.log.Printf("{\"ts\":\"%s\",\"level\":\"error\",\"msg\":\"marshal_failed\"}", payload.Timestamp)
		return
	}
	global.log.Println(string(data))
}

func Info(msg string, fields map[string]interface{}) {
	logEntry("info", msg, fields)
}

func Warn(msg string, fields map[string]interface{}) {
	logEntry("warn", msg, fields)
}

func Error(msg string, fields map[string]interface{}) {
	logEntry("error", msg, fields)
}