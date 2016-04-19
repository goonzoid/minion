package main

import (
	"os"
	"time"
)

func main() {
	os.Stdout.WriteString("hello stdout\n")
	os.Stderr.WriteString("hello stderr\n")

	for {
		time.Sleep(time.Minute)
	}
}
