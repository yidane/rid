package log

import (
	syslog "log"
)

//Error display error
func Error(v ...interface{}) {
	syslog.Println(v)
}

//Succeed display successful information
func Succeed(v ...interface{}) {
	syslog.Println(v)
}

//Warn display warn
func Warn(v ...interface{}) {
	syslog.Println(v)
}

//Info display information
func Info(v ...interface{}) {
	syslog.Println(v)
}
