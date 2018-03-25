package main

import "fmt"

// type Message func(string)

// func Intro(s string) {

// }
func RenderMessageHeading(msg string) {

	fmt.Println(" ----------------------------------------------------------------- ")
	fmt.Println("|                       " + msg + "                          |")
	fmt.Println(" ----------------------------------------------------------------- ")
	fmt.Println("\n\n")
}

func RenderMessageEnd() {
	fmt.Println("\n\n")
	fmt.Print("-------------------------------------------------------------------")
	fmt.Println("\n\n")
}

func Message(msg string) {
	fmt.Println("> " + msg)
}
