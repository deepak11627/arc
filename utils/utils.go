package utils

import "fmt"

func RenderMessageHeading(msg string) {

	fmt.Println("-----------------------------------------------------------------")
	fmt.Println(msg)
	fmt.Println("-----------------------------------------------------------------")
	fmt.Println("\n")
}

func RenderMessageEnd() {
	fmt.Println("\n")
	fmt.Print("-----------------------------------------------------------------")
	fmt.Println("\n")
}

func Message(msg string) {
	fmt.Println("> " + msg)
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
