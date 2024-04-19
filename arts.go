package main

import "fmt"

// const Start = ` __        _______ _     ____ ___  __  __ _____   _____ ___    ____  ____  __  __
//  \ \      / / ____| |   / ___/ _ \|  \/  | ____| |_   _/ _ \  / ___||  _ \|  \/  |
//   \ \ /\ / /|  _| | |  | |  | | | | |\/| |  _|     | || | | | \___ \| | | | |\/| |
//    \ V  V / | |___| |__| |__| |_| | |  | | |___    | || |_| |  ___) | |_| | |  | |
//     \_/\_/  |_____|_____\____\___/|_|  |_|_____|   |_| \___/  |____/|____/|_|  |_|
// `

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorGreen  Color = "\u001b[32m"
	ColorRed    Color = "\u001b[31m"
	ColorYellow Color = "\u001b[33m"
	ColorBlue   Color = "\u001b[34m"
	ColorReset  Color = "\u001b[0m"
)

func Colorize(color Color, message string) string {
	return fmt.Sprintf("%v%v%v", string(color), message, string(ColorReset))
}

func LogsColorize(color Color, logType string, args ...string) string {
	outs := string(color) + logType
	for _, str := range args {
		outs += " " + str
	}
	outs += string(ColorReset)
	return outs
}
