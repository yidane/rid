package util

import "fmt"

//Green change the output's color in terminal
func Green(v ...interface{}) string {
	return fmt.Sprintf("%c[1;0;32m%s%c[0m", 0x1B, v, 0x1B)
}

//Red change the output's color in terminal
func Red(v ...interface{}) string {
	return fmt.Sprintf("%c[1;0;31m%s%c[0m", 0x1B, v, 0x1B)
}

//Yellow change the output's color in terminal
func Yellow(v ...interface{}) string {
	return fmt.Sprintf("%c[1;0;33m%s%c[0m", 0x1B, v, 0x1B)
}
