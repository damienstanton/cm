package main

//#include <stdio.h>
import "C"

// HiFromGo exports a string message
//export HiFromGo
func HiFromGo() *C.char {
	return C.CString("Hi from Go!")
}

func main() {}
