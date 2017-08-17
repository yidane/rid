package log

import (
	syslog "log"
)

func init(){
	//判断操作系统，分别做日志颜色处理
}

//Error display error
func Error(v ...interface{}) {
	syslog.Printf("%c[1;0;31m%s%c[0m\n", 0x1B, v, 0x1B)
}

//Succeed display successful infomation
func Succeed(v ...interface{}) {
	syslog.Printf("%c[1;0;32m%s%c[0m\n", 0x1B, v, 0x1B)
}

//Warn display warn
func Warn(v ...interface{}) {
	syslog.Printf("%c[1;0;33m%s%c[0m\n", 0x1B, v, 0x1B) 
}

//Info display information
func Info(v ...interface{}) {
	syslog.Println(v)
}
