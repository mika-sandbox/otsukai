package otsukai

import (
	"fmt"
	"github.com/fatih/color"
)

const HEADER_STR = "[OTSUKAI]"

var blue = color.New(color.FgHiBlue).SprintfFunc()
var green = color.New(color.FgHiGreen).SprintfFunc()
var red = color.New(color.FgHiRed).SprintfFunc()
var yellow = color.New(color.FgHiYellow).SprintfFunc()

func Debug() {}

func Info(msg string, extras ...any) {
	fmt.Printf("%s %s\n", blue(HEADER_STR+" [INFO]    "), fmt.Sprintf(msg, extras...))
}

func Success(msg string, extras ...any) {
	fmt.Printf("%s %s\n", green(HEADER_STR+" [SUCCESS] "), fmt.Sprintf(msg, extras...))
}

func Warn(msg string, extras ...any) {
	fmt.Printf("%s %s\n", yellow(HEADER_STR+" [WARNING] "), fmt.Sprintf(msg, extras...))
}

func Err(msg string, extras ...any) {
	fmt.Printf("%s %s\n", red(HEADER_STR+" [ERROR]   "), fmt.Sprintf(msg, extras...))
}
