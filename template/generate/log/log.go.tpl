package log

import (
	"bytes"
	"{{.packagePrefix}}/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type (
	ZapLog struct {
		zapLog       *zap.SugaredLogger
		currLogLevel zap.AtomicLevel
	}
)

var (
	DefaultLog *ZapLog
)

func init() {
	logLevelStr := config.GetValue("log", "level")
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logLevel := zap.NewAtomicLevel()
	err := logLevel.UnmarshalText([]byte(logLevelStr))
	if err != nil {
		log.Fatalf("logLevel UnmarshalText failed, error = %s!", err)
	}

	programName := strings.Split(filepath.Base(os.Args[0]), ".")[0]

	workPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("get workPath failed, error = %s!", err)
	}

	workPath, err = filepath.Abs(workPath)
	if err != nil {
		log.Fatalf("get workPath Abs failed, error = %s!", err)
	}

	logPath := workPath + "/logs"
	err = os.MkdirAll(logPath, 0777)
	if err != nil {
		log.Fatalf("create log path failed, error = %s!", err)
	}

	var filePathNames = []string{"./logs/" + programName + ".log"}
	if logLevelStr == "debug" {
		filePathNames = append(filePathNames, "stdout")
	}

	customCfg := zap.Config{
		Level:            logLevel,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      filePathNames,
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := customCfg.Build()
	if err != nil {
		log.Fatalf("customCfg Build failed, error = %s!", err)
	}

	newLogger := logger.Named(programName).WithOptions(zap.AddCallerSkip(1))
	sugar := newLogger.Sugar()

	DefaultLog = &ZapLog{
		zapLog:       sugar,
		currLogLevel: logLevel,
	}
}

func (zapLog *ZapLog) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	zapLog.zapLog.Error(string(p))
	return len(p), nil
}

func GetLogger() *ZapLog {
	return DefaultLog
}

func Sync() {
	err := DefaultLog.zapLog.Sync()
	if err != nil {
		log.Fatalf("Sync failed, error = %s!", err)
	}
}

func Debug(msg string, keysAndValues ...interface{}) {
	if !judgeDefaultLogger() {
		log.Println("DEBUG:", msg, keysAndValues)
		return
	}
	DefaultLog.zapLog.Debugw(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	if !judgeDefaultLogger() {
		log.Println("INFO:", msg, keysAndValues)
		return
	}
	DefaultLog.zapLog.Infow(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	if !judgeDefaultLogger() {
		log.Println("WARN:", msg, keysAndValues)
		return
	}
	DefaultLog.zapLog.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	if !judgeDefaultLogger() {
		log.Println("ERROR:", msg, keysAndValues)
		return
	}
	DefaultLog.zapLog.Errorw(msg, keysAndValues...)
}

func SetLogLevel(logLevel string) error {
	return DefaultLog.currLogLevel.UnmarshalText([]byte(logLevel))
}

func GetLogLevel() string {
	return DefaultLog.currLogLevel.String()
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}

func judgeDefaultLogger() bool {
	if DefaultLog == nil {
		return false
	} else if DefaultLog.zapLog == nil {
		return false
	}
	return true
}
