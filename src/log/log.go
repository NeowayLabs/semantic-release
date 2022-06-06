package log

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	errorLevel = "error"
	warnLevel  = "warn"
	fatalLevel = "fatal"
)

type Log struct {
	program       logProgram
	level         string
	log           log.Logger
	exitWhenFatal bool
}

type logProgram struct {
	name    string
	version string
}

func (l *Log) parseLevel() (string, error) {
	switch strings.ToLower(l.level) {
	case "fatal":
		return fatalLevel, nil
	case "error":
		return errorLevel, nil
	case "warn", "warning":
		return warnLevel, nil
	case "info":
		return infoLevel, nil
	case "debug":
		return debugLevel, nil
	}

	return "", errors.New("not a valid log Level: " + l.level)
}

func (l *Log) format(message string) string {
	return fmt.Sprintf("%s - %s:", l.program.name, l.program.version) + message
}

func (l *Log) Info(s string, args ...interface{}) {
	l.log.Println(l.format(fmt.Sprintf(s, args...)))
}

func (l *Log) Debug(s string, args ...interface{}) {
	if l.level == debugLevel {
		l.log.Println(l.format(fmt.Sprintf(s, args...)))
	}
}

func (l *Log) Error(s string, args ...interface{}) {
	if l.level == debugLevel || l.level == errorLevel {
		l.log.Println(l.format(fmt.Sprintf(s, args...)))
	}
}

func (l *Log) Warn(s string, args ...interface{}) {
	if l.level == debugLevel || l.level == warnLevel {
		l.log.Println(l.format(fmt.Sprintf(s, args...)))
	}
}

func (l *Log) Fatal(s string, args ...interface{}) {
	if l.level == debugLevel || l.level == fatalLevel {
		l.log.Printf(l.format(fmt.Sprintf(s, args...)))

		if l.exitWhenFatal {
			os.Exit(1)
		}
	}
}

func (l *Log) Configure(output *bytes.Buffer, notExitWhenFatal bool) {
	if output != nil {
		l.log.SetOutput(output)
	}

	if notExitWhenFatal {
		l.exitWhenFatal = false
	}
}

func New(name string, version string, level string) (*Log, error) {

	newLog := &Log{
		program:       logProgram{name: name, version: version},
		level:         level,
		log:           *log.Default(),
		exitWhenFatal: true,
	}

	// fmt.Println("hi")
	logLevel, err := newLog.parseLevel()
	if err != nil {
		return nil, err
	}

	newLog.level = logLevel

	return newLog, nil
}
