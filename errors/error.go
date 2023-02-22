package errors

import (
	"fmt"
	"os"
	"syscall"
)

const (
	COLOR_RESET string = "\033[0m"
	RED         string = "\033[31m"
	GREEN       string = "\033[32m"
	YELLOW      string = "\033[33m"
	BLUE        string = "\033[34m"
	PURPLE      string = "\033[35m"
	CYAN        string = "\033[36m"
	WHITE       string = "\033[37m"
)

var FUCK_YOU bool = false
var WALL_FLAG bool = false

func ThrowWarning(text string, quit bool) {
	fmt.Println(PURPLE, "\r[WARNING]", text, COLOR_RESET)
	if quit {
		os.Exit(0)
	}
}

func ThrowError(text string, quit bool) {
	fmt.Println(RED, "\n\r[ERROR]", text, COLOR_RESET)
	if FUCK_YOU {
		// hard crash
		syscall.Kill(1, syscall.SIGSEGV)
	}
	if quit {
		os.Exit(0)
	}

}

func ThrowInfo(text string, quit bool) {
	fmt.Println(BLUE, "\r[INFO]", text, COLOR_RESET)
	if quit {
		os.Exit(0)
	}
}

func ThrowSuccess(text string, quit bool) {
	fmt.Println(GREEN, "\r[SUCCESS]", text, COLOR_RESET)
	if quit {
		os.Exit(0)
	}
}

func ThrowCustom(color, text string, quit bool) {
	fmt.Println(color, "\r", text, COLOR_RESET)
	if quit {
		os.Exit(0)
	}
}
