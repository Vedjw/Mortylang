package main

import (
	"bufio"
	"fmt"
	"morty/repl"
	"os"
	"os/user"
	"regexp"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s! This is the Morty programming language!\n", user.Name)
	fmt.Printf("Type name of the file to execute\n")

	var input string
	fmt.Scanln(&input)
	if !isMorty(input) {
		fmt.Printf("%s is not a .morty file ", input)
		os.Exit(1)
	}

	filename := input

	inFile, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Error opening the %s ", filename)
		os.Exit(1)
	}
	defer inFile.Close()
	r := bufio.NewReader(inFile)

	repl.Start(r, os.Stdout)
}

func isMorty(name string) bool {
	matched, err := regexp.MatchString(".morty", name)
	if err != nil {
		fmt.Print(err)
	}

	return matched
}
