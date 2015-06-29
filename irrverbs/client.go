package main

import (
	"fmt"
	"os"
)

func main() {
	if !checkArgs() {
		usage()
		os.Exit(1)
	}
}

func checkArgs() bool {
	if len(os.Args) < 2 {
		return false
	} else if os.Args[1] != "add" && os.Args[1] != "learn" {
		return false
	} else if os.Args[1] == "add" && len(os.Args) != 5 {
		return false
	} else if os.Args[1] == "learn" && len(os.Args) != 2 {
		return false
	}
	return true
}

func usage() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s add v1 v2 v3\n", os.Args[0])
	fmt.Printf("\t%s learn\n", os.Args[0])
}
