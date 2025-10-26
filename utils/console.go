package utils

import "fmt"

func ClearScreen() {
	fmt.Print("\033[H\033[2J") // ANSI escape code to clear screen
}
