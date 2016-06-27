package log

import "strings"

var log logger

func init() {
	log = newBeegoAdapter()
}

type Level int

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelError
	LevelFatal
)

type Adapter struct {
	Config string
}

type LogConfig struct {
	Level    Level
	Adapters map[string]Adapter
}

//Log  日志接口
type logger interface {
	SetLog(*LogConfig) error
	//Trace 跟踪
	Trace(format string, args ...interface{})

	//Debug 调试
	Debug(format string, args ...interface{})

	//Info 信息
	Info(format string, args ...interface{})

	//Error 错误
	Error(format string, args ...interface{})

	//Fatal 致命错误
	Fatal(format string, args ...interface{})
}

//SetLog 配置log
func SetLog(conf *LogConfig) error {
	return log.SetLog(conf)
}

//Tracef 跟踪
func Tracef(format string, args ...interface{}) {
	log.Trace(format, args)
}

//Debugf 调试
func Debugf(format string, args ...interface{}) {
	log.Debug(format, args)
}

//Infof 信息
func Infof(format string, args ...interface{}) {
	log.Info(format, args)
}

//Errorf 错误
func Errorf(format string, args ...interface{}) {
	log.Error(format, args)
}

//Fatalf 致命错误
func Fatalf(format string, args ...interface{}) {
	log.Fatal(format, args)
}

//Trace 跟踪
func Trace(v ...interface{}) {
	log.Trace(generateFmtStr(len(v)), v...)
}

//Debug 调试
func Debug(v ...interface{}) {
	log.Debug(generateFmtStr(len(v)), v...)
}

//Info 信息
func Info(v ...interface{}) {
	log.Info(generateFmtStr(len(v)), v...)
}

//Error 错误
func Error(v ...interface{}) {
	log.Error(generateFmtStr(len(v)), v...)
}

//Fatal 致命错误
func Fatal(v ...interface{}) {
	log.Fatal(generateFmtStr(len(v)), v...)
}

func generateFmtStr(n int) string {
	return strings.Repeat("%v ", n)
}
