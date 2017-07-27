package log

import syslog "log"

//Error display error
func Error(i ...interface{}) {
	syslog.Println(i)
}

//Succeed display successful infomation
func Succeed(i ...interface{}) {

}

//Warn display warn
func Warn(i ...interface{}) {
	syslog.Println(i)
}

//Info display information
func Info(i ...interface{}) {
	syslog.Println(i)
}
