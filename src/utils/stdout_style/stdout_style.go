package stdoutstyle

import "runtime"

var (
	Reset    = "\033[0m"
	Red      = "\033[31m"
	Green    = "\033[32m"
	Yellow   = "\033[33m"
	Blue     = "\033[34m"
	Purple   = "\033[35m"
	Cyan     = "\033[36m"
	Gray     = "\033[37m"
	White    = "\033[97m"
	BGBlack  = "\033[40;1;37m"
	BGRed    = "\033[41;1;37m"
	BGGreen  = "\033[42;1;37m"
	BGBlue   = "\033[44;1;37m"
	BGPurple = "\033[45;1;37m"
	BGCyan   = "\033[46;1;37m"
	BGGray   = "\033[47;1;37m"
	BGWhite  = "\033[47m"
	Blink    = "\033[5;30m"
)

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
		BGBlack = ""
		BGRed = ""
		BGGreen = ""
		BGBlue = ""
		BGPurple = ""
		BGCyan = ""
		BGGray = ""
		Blink = ""
	}
}
