package utils

import "fmt"

func RenderMessageHeading(msg string) {
	fmt.Print("-----------------------------------------------------------------\n" + msg)
}

func RenderMessageEnd() {
	fmt.Println("\n-----------------------------------------------------------------")
}

func Message(msg string) {
	fmt.Print("\n> " + msg)
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
