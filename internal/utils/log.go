package utils

import (
	"fmt"
	"runtime"
)

func PrepareLogMsg(msg string) string {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	return fmt.Sprintf("Error: %s, Function Name: %s, Line: %d, File: %s", msg, funcName, line, file)
}
