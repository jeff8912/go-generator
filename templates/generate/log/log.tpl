package log

import (
	"bytes"
	"{{.packagePrefix}}/common"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

var defaultLog *Zlog

type Zlog struct {
	zapLog       *zap.SugaredLogger
	currLogLevel zap.AtomicLevel
}

// this function for echo
func (zlog *Zlog) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	zlog.zapLog.Error(string(p))
	return len(p), nil
}

// this function for zk lib
func (zlog *Zlog) Printf(template string, args ...interface{}) {
	zlog.zapLog.Errorf("zk_lib "+template, args)
}

func Init(programName, logLevelStr string, filenames ...string) *Zlog {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "T",
		LevelKey:   "L",
		NameKey:    "N",
		CallerKey:  "C",
		MessageKey: "M",
		//StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	logLevel := zap.NewAtomicLevel()
	logLevel.UnmarshalText([]byte(logLevelStr))

	logPath := common.JoinPath(common.GetMainPath(), "log")
	err := os.MkdirAll(logPath, 0777)
	if err != nil {
		log.Fatalf("create log path failed, error = %s!", err)
	}

	filePathNames := append(filenames, common.JoinPath(logPath, programName+".log"))

	customCfg := zap.Config{
		Level:            logLevel,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    encoderConfig,
		OutputPaths:      filePathNames,
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := customCfg.Build()
	new_logger := logger.Named(programName).WithOptions(zap.AddCallerSkip(1))
	sugar := new_logger.Sugar()

	defaultLog = &Zlog{
		zapLog:       sugar,
		currLogLevel: logLevel,
	}

	return defaultLog
}

func GetLogger() *Zlog {
	return defaultLog
}

func Sync() {
	defaultLog.zapLog.Sync()
}

func Debug(msg string, keysAndValues ...interface{}) {
	if !judgeDefaultLogger() {
		log.Println("DEBUG:", msg, keysAndValues)
		return
	}
	defaultLog.zapLog.Debugw(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...interface{}) {
	if !judgeDefaultLogger() {
		log.Println("INFO:", msg, keysAndValues)
		return
	}
	defaultLog.zapLog.Infow(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	if !judgeDefaultLogger() {
		log.Println("WARN:", msg, keysAndValues)
		return
	}
	defaultLog.zapLog.Warnw(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	if !judgeDefaultLogger() {
		log.Println("ERROR:", msg, keysAndValues)
		return
	}
	defaultLog.zapLog.Errorw(msg, keysAndValues...)
}

func SetLogLevel(logLevel string) {
	defaultLog.currLogLevel.UnmarshalText([]byte(logLevel))
}

func GetLogLevel() string {
	return defaultLog.currLogLevel.String()
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}

func judgeDefaultLogger() bool {
	if defaultLog == nil {
		return false
	} else if defaultLog.zapLog == nil {
		return false
	}
	return true
}
