package log

import (
	"testing"
)

func TestGetLogger(t *testing.T) {
	t.Log(GetLogger().currLogLevel)
}

func TestDebug(t *testing.T) {
	Debug("debug example", "key", "value")
}

func TestInfo(t *testing.T) {
	Info("info example", "key", "value")
}

func TestWarn(t *testing.T) {
	Warn("warn example", "key", "value")
}

func TestError(t *testing.T) {
	Error("error example", "key", "value")
}

func TestSetLogLevel(t *testing.T) {
	err := SetLogLevel("error")
	if err != nil {
		panic(err)
	}

	Debug("error example", "key", "value")
	Info("error example", "key", "value")
	Warn("error example", "key", "value")
	Error("error example", "key", "value")
}

func TestGetLogLevel(t *testing.T) {
	t.Log(GetLogLevel())
}
