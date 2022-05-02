//go:build unit
// +build unit

package log_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/NeowayLabs/semantic-release/src/log"
	"github.com/NeowayLabs/semantic-release/src/tests"
)

var (
	output bytes.Buffer
)

func NewLogger(outputLog bytes.Buffer, name, version, level string) (*log.Log, *bytes.Buffer) {
	logger, err := log.New(name, version, level)
	if err != nil {
		return nil, nil
	}
	logger.Configure(&outputLog, true)
	return logger, &outputLog
}

func TestNewErrorLevel(t *testing.T) {
	_, err := log.New("test", "0.1", "not-found")
	tests.AssertError(t, err)
}

func TestInfo(t *testing.T) {
	// var output bytes.Buffer
	logger, outputLog := NewLogger(output, "test", "1.0.0", "info")
	tests.AssertNotNil(t, logger)
	logger.Info("Hi")

	splited := strings.Split(outputLog.String(), " ")
	actual := strings.TrimSpace(splited[2]) + splited[3] + strings.TrimSpace(splited[4])
	tests.AssertEqualValues(t, "test-1.0.0:Hi", actual)
}

func TestDebug(t *testing.T) {
	logger, outputLog := NewLogger(output, "test", "1.0.0", "debug")
	tests.AssertNotNil(t, logger)
	logger.Debug("Hi")

	splited := strings.Split(outputLog.String(), " ")
	actual := strings.TrimSpace(splited[2]) + splited[3] + strings.TrimSpace(splited[4])
	tests.AssertEqualValues(t, "test-1.0.0:Hi", actual)
}

func TestError(t *testing.T) {
	logger, outputLog := NewLogger(output, "test", "1.0.0", "error")
	tests.AssertNotNil(t, logger)
	logger.Error("Hi")

	splited := strings.Split(outputLog.String(), " ")
	actual := strings.TrimSpace(splited[2]) + splited[3] + strings.TrimSpace(splited[4])
	tests.AssertEqualValues(t, "test-1.0.0:Hi", actual)
}

func TestWarn(t *testing.T) {
	logger, outputLog := NewLogger(output, "test", "1.0.0", "warn")
	tests.AssertNotNil(t, logger)
	logger.Warn("Hi")

	splited := strings.Split(outputLog.String(), " ")
	actual := strings.TrimSpace(splited[2]) + splited[3] + strings.TrimSpace(splited[4])
	tests.AssertEqualValues(t, "test-1.0.0:Hi", actual)
}

func TestFatal(t *testing.T) {
	logger, outputLog := NewLogger(output, "test", "1.0.0", "fatal")
	tests.AssertNotNil(t, logger)
	logger.Fatal("Hi")

	splited := strings.Split(outputLog.String(), " ")
	actual := strings.TrimSpace(splited[2]) + splited[3] + strings.TrimSpace(splited[4])
	tests.AssertEqualValues(t, "test-1.0.0:Hi", actual)
}
