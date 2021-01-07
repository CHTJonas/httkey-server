package main

import (
	"fmt"
	"os"
)

func main() {
	rawurl := os.Args[1]
	hash, err := RawURLToHash(rawurl)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Printf("%s: %s\n", rawurl, hash)
}
