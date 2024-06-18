package main

import (
	"fmt"
	"golang.org/x/sys/unix"
)

func recoveryShell() {
	unix.Exec("/bin/sh", []string{"sh"}, []string{})
}

func main() {
	fmt.Println("Hello, world!")
	fmt.Println("Welcum to init uwu")

	recoveryShell()
}
