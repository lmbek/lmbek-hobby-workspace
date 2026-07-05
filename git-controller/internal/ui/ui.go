package ui

import (
	"fmt"
	"strings"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

func Header(text string) {
	fmt.Printf("\n%s%s%s\n%s\n", ColorBold, ColorCyan, strings.ToUpper(text), strings.Repeat("=", len(text)))
}

func Info(format string, a ...any) {
	fmt.Printf("%s%s%s\n", ColorBlue, fmt.Sprintf(format, a...), ColorReset)
}

func Success(format string, a ...any) {
	fmt.Printf("%s✔ %s%s\n", ColorGreen, fmt.Sprintf(format, a...), ColorReset)
}

func Warn(format string, a ...any) {
	fmt.Printf("%s⚠ %s%s\n", ColorYellow, fmt.Sprintf(format, a...), ColorReset)
}

func Error(format string, a ...any) {
	fmt.Printf("%s✖ %s%s\n", ColorRed, fmt.Sprintf(format, a...), ColorReset)
}

func Step(num int, text string) {
	fmt.Printf("\n%s[%d] %s%s\n", ColorBold, num, text, ColorReset)
}

func Note(text string) {
	fmt.Printf("\n%sNote: %s%s\n", ColorCyan, text, ColorReset)
}
