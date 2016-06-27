package log

import "github.com/astaxie/beego/logs"

type beegoAdapter struct {
	beegoLog *logs.BeeLogger
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
