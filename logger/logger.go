package logger

import (
	"fmt"
	"time"
)

// Infor .
func Infor(format string, v ...interface{}) {
	print(36, format, v...)
}

// Debug .
func Debug(format string, v ...interface{}) {
	print(37, format, v...)
}

// Error .
func Error(format string, v ...interface{}) {
	print(31, format, v...)
}

// print .
func print(level int, format string, v ...interface{}) {
	fmt.Print(fmt.Sprintf("[%s]--------%c[%d;%d;%dm%s%c[0m\n", time.Now().Format("2006-01-02 15:04:05"), 0x1B, 0 /*字体*/, 0 /*背景*/, level /*前景*/, fmt.Sprintf(format, v...), 0x1B))
}
