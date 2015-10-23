package logger

import (
	"fmt"
	"log"
	"os"
	"path"
)

const (
	_ = iota
	levelError
	levelWarning
	levelNotice
	levelInfo
	levelDebug

	defaultLevel string = "info"
)

var (
	logger   *log.Logger
	logLevel int

	levelMap = map[string]int{
		"error":   levelError,
		"warning": levelWarning,
		"notice":  levelNotice,
		"info":    levelInfo,
		"debug":   levelDebug,
	}
)

// Init initializes the logging messages handler
func Init(logPath, level string) error {
	var (
		logOut *os.File
		err    error
	)

	if level == "" {
		level = defaultLevel
	}

	if logPath != "" && logPath != "-" {
		dirPath, _ := path.Split(logPath)

		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}

		logOut, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("unable to open log file: %s", err)
		}
		defer logOut.Close()
	} else {
		logOut = os.Stderr
	}

	logger = log.New(logOut, "", log.LstdFlags|log.Lmicroseconds)

	logLevel, err = getLevelByName(level)
	if err != nil {
		return err
	}

	return nil
}

// Error prints an error logging message
func Error(context, format string, v ...interface{}) {
	printLog(levelError, context, format, v...)
}

// Warning prints a warning logging message
func Warning(context, format string, v ...interface{}) {
	printLog(levelWarning, context, format, v...)
}

// Notice prints a notice logging message
func Notice(context, format string, v ...interface{}) {
	printLog(levelNotice, context, format, v...)
}

// Info prints an info logging message
func Info(context, format string, v ...interface{}) {
	printLog(levelInfo, context, format, v...)
}

// Debug prints a debug logging message
func Debug(context, format string, v ...interface{}) {
	printLog(levelDebug, context, format, v...)
}

func getLevelByName(name string) (int, error) {
	level, ok := levelMap[name]
	if !ok {
		return 0, fmt.Errorf("invalid level `%s'", name)
	}

	return level, nil
}

func printLog(level int, context, format string, v ...interface{}) {
	var crit string

	if level > logLevel {
		return
	}

	switch level {
	case levelError:
		crit = "ERROR"

	case levelWarning:
		crit = "WARNING"

	case levelNotice:
		crit = "NOTICE"

	case levelInfo:
		crit = "INFO"

	case levelDebug:
		crit = "DEBUG"
	}

	logger.Printf("%s: %s: %s", crit, context, fmt.Sprintf(format, v...))
}
