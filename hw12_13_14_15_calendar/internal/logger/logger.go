package logger

import (
	"log"
	"strings"
)

type Logger struct {
	level string
}

func New(level string) *Logger {
	return &Logger{
		level: strings.ToLower(level),
	}
}

func (l Logger) logMessage(level, msg string) {
	levels := map[string]int{
		"debug": 4,
		"info":  3,
		"warn":  2,
		"error": 1,
	}

	currentLevel, exists := levels[l.level]
	msgLevel, valid := levels[level]

	if !exists || !valid || currentLevel < msgLevel {
		return // no need to log then
	}

	log.Printf("[%s] %s\n", strings.ToUpper(level), msg)
}

func (l Logger) Info(msg string) {
	l.logMessage("info", msg)
}

func (l Logger) Error(msg string) {
	l.logMessage("error", msg)
}

func (l Logger) Debug(msg string) {
	l.logMessage("debug", msg)
}

func (l Logger) Warn(msg string) {
	l.logMessage("warn", msg)
}
