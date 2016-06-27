package log

import "github.com/astaxie/beego/logs"

type beegoAdapter struct {
	beegoLog *logs.BeeLogger
}

func (log *beegoAdapter) SetLog(conf *LogConfig) error {
	switch conf.Level {
	case LevelTrace:
		log.beegoLog.SetLevel(logs.LevelTrace)
	case LevelDebug:
		log.beegoLog.SetLevel(logs.LevelDebug)
	case LevelInfo:
		log.beegoLog.SetLevel(logs.LevelInfo)
	case LevelFatal:
		log.beegoLog.SetLevel(logs.LevelEmergency)
	}

	for k := range conf.Adapters {
		if err := log.beegoLog.SetLogger(k, conf.Adapters[k].Config); err != nil {
			return err
		}
	}
	return nil
}

//Trace 跟踪
func (log *beegoAdapter) Trace(format string, args ...interface{}) {
	log.beegoLog.Trace(format, args)
}

//Debug 调试
func (log *beegoAdapter) Debug(format string, args ...interface{}) {
	log.beegoLog.Debug(format, args)
}

//Info 信息
func (log *beegoAdapter) Info(format string, args ...interface{}) {
	log.beegoLog.Info(format, args)
}

//Error 错误
func (log *beegoAdapter) Error(format string, args ...interface{}) {
	log.beegoLog.Error(format, args)
}

//Fatal 致命错误
func (log *beegoAdapter) Fatal(format string, args ...interface{}) {
	log.beegoLog.Emergency(format, args)
}

func newBeegoAdapter() *beegoAdapter {
	var ret = &beegoAdapter{beegoLog: logs.NewLogger(100)}
	return ret
}
