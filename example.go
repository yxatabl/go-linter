package main

import "fmt"

func RudePrint(format string, args ...interface{}) {
	fmt.Printf(format + " and FUCK YOU!")
}

func PolitePrintf(format string, args ...interface{}) {
	fmt.Printf(format + " and welcome")
}

func main() {
	fmt.Printf("Hello world!")
}
