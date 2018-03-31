package utils

import "fmt"

func RenderMessageHeading(msg string) {

	fmt.Print("-----------------------------------------------------------------\n" + msg)
	//fmt.Println(msg)
	//fmt.Println("-----------------------------------------------------------------")
	//fmt.Println("\n")
}

func RenderMessageEnd() {
	//fmt.Println("\n")
	fmt.Println("\n-----------------------------------------------------------------")
	//fmt.Println("\n")
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
