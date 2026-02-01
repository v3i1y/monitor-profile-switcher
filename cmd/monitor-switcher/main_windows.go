//go:build windows

package main

import (
	"fmt"
	"os"
	"strings"

	"monitor-profile-switcher/internal/switcher"
)

type command struct {
	kind  string
	value string
}

func main() {
	args := os.Args[1:]
	debug := false
	noIDMatch := false
	virtualInject := false
	var commands []command

	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") {
			commands = append(commands, command{kind: "load", value: arg})
			continue
		}
		key, value := splitArg(arg)
		switch strings.ToLower(key) {
		case "-debug":
			debug = true
			fmt.Println("Debug output enabled")
		case "-noidmatch":
			noIDMatch = true
			if debug {
				fmt.Println("Disabled matching of adapter IDs")
			}
		case "-v":
			virtualInject = true
			if debug {
				fmt.Println("Enabled virtual desktop injection")
			}
		case "-save":
			commands = append(commands, command{kind: "save", value: value})
		case "-load":
			commands = append(commands, command{kind: "load", value: value})
		case "-print":
			commands = append(commands, command{kind: "print"})
		}
	}

	if len(commands) == 0 {
		printUsage()
		return
	}

	for _, cmd := range commands {
		switch cmd.kind {
		case "save":
			path, err := switcher.ResolveProfilePath(cmd.value, true)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid -save argument:", err)
				os.Exit(1)
			}
			if err := switcher.SaveProfile(path, debug); err != nil {
				fmt.Fprintln(os.Stderr, "Save failed:", err)
				os.Exit(1)
			}
		case "load":
			path, err := switcher.ResolveProfilePath(cmd.value, false)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Invalid -load argument:", err)
				os.Exit(1)
			}
			if err := switcher.LoadProfile(path, debug, noIDMatch, virtualInject); err != nil {
				fmt.Fprintln(os.Stderr, "Load failed:", err)
				os.Exit(1)
			}
		case "print":
			if err := switcher.PrintSummary(os.Stdout); err != nil {
				fmt.Fprintln(os.Stderr, "Print failed:", err)
				os.Exit(1)
			}
		}
	}
}

func splitArg(arg string) (string, string) {
	parts := strings.SplitN(arg, ":", 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func printUsage() {
	fmt.Println("Monitor Profile Switcher (Go CLI)\n")
	fmt.Println("Parameters:")
	fmt.Println("  -save:{file}        save current monitor configuration to file")
	fmt.Println("  -load:{file}        load and apply monitor configuration from file")
	fmt.Println("  -debug              enable debug output (use before -load or -save)")
	fmt.Println("  -noidmatch          disable matching of adapter IDs")
	fmt.Println("  -v                  enable virtual desktop injection (advanced)")
	fmt.Println("  -print              print current monitor configuration summary")
	fmt.Println("")
	fmt.Println("If {file} is a filename (no path), it is stored under:")
	fmt.Println("  %USERPROFILE%\\Monitor Profiles")
	fmt.Println("If {file} has no extension, .monitorprofile is added.")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  monitor-switcher.exe -save:Profile.json")
	fmt.Println("  monitor-switcher.exe -load:Profile.json")
	fmt.Println("  monitor-switcher.exe -debug -load:Profile.json")
}
