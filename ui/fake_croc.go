//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// Fake croc that prints Code is: fakecode, then progress, then exits
	fmt.Println("Code is: fakecode")
	time.Sleep(50 * time.Millisecond)
	for i := 1; i <= 100; i += 25 {
		fmt.Printf("Sending 'file.txt' %d%%\n", i)
		time.Sleep(30 * time.Millisecond)
	}
	fmt.Println("Sending complete")
	os.Exit(0)
}
