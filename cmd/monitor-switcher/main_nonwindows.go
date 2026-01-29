//go:build !windows

package main

import "fmt"

func main() {
	fmt.Println("monitor-switcher is supported on Windows only.")
}
