package internal

import (
	"fmt"

	"github.com/fatih/color"
)

func ErrPrinter(msg error) {
	color.Red(fmt.Sprintln("Error " +msg.Error()))
}

func CrrPrinter(msg string) {
	color.Green(msg)
}

func WarnPrinter(msg string) {
	color.Yellow(msg)
}

func ErrColorString(msg string) string {
	return color.RedString(msg)
}

func CrrColorString(msg string) string {
	return color.GreenString(msg)
}