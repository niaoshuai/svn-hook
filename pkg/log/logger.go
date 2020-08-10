/**
 * @Author: niaoshuai
 * @Date: 2020/8/9 8:51 上午
 */
package log

import (
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

var (
	Logger = log.New()
)

func InitLog(logPath string) {
	//打开日志文件
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("open log file err, file path:", logPath, ",err:", err)
	}

	//修改日志的输出方式
	Logger.Out = file
	//设置日志格式为json
	Logger.Formatter = &log.JSONFormatter{}
}

func Info(msg string) {
	fields := formatLog()
	Logger.WithFields(fields).Info(msg)
}

func Fatal(err error) {
	fields := formatLog()
	Logger.WithFields(fields).Fatal(err)
}

func Error(err error) {
	fields := formatLog()
	Logger.WithFields(fields).Error(err)
}

/**
 * 为日志字段增加文件和行号
 */
func formatLog() log.Fields {
	_, file, line, ok := runtime.Caller(2)

	if !ok {
		log.Fatalln("runtime caller err")
	}

	var fields = log.Fields{
		"file": file,
		"line": line,
	}

	return fields
}
