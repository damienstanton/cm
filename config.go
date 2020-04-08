package main

import (
	"fmt"
	"time"
)

const catchVersion = "v2.11.3"
const cmVersion = "v0.1.0"

var compileTimeout = 10 * time.Minute

// printBanner shows the current cm version
func printBanner() {
	fmt.Println("╔═════════════════════════╗")
	fmt.Printf("║ Compiler Manager %s ║\n", cmVersion)
	fmt.Println("╚═════════════════════════╝")
}

// TODO: JSON parser
