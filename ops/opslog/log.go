package opslog

import (
	"log"
	"os"
)

/**
 * @Author: yangchujie
 * @Author: chujieyang@gmail.com
 * @Date: 6/16/21 10:54 AM
 * @Desc:
 */

var logInfo = log.New(os.Stdout, " <INFO> ", log.Lshortfile|log.Ldate|log.Ltime)
var logWarn = log.New(os.Stdout, " <WARN> ", log.Lshortfile|log.Ldate|log.Ltime)
var logError = log.New(os.Stdout, " <ERROR> ", log.Lshortfile|log.Ldate|log.Ltime)

func Info() *log.Logger {
	return logInfo
}

func Warn() *log.Logger {
	return logWarn
}

func Error() *log.Logger {
	return logError
}