package main

import "fmt"

func green(s string) string {
	return fmt.Sprintf("%c[1;0;32m%s%c[0m", 0x1B, s, 0x1B)
}

func red(s string) string {
	return fmt.Sprintf("%c[1;0;32m%s%c[0m", 0x1B, s, 0x1B)
}
