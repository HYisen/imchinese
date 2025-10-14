package main

import (
	"fmt"
	"imchinese/finder"
	"log"
	"os"
)

func main() {
	data, err := os.ReadFile("text.md")
	if err != nil {
		log.Fatal(err)
	}
	prettyPrint(finder.Find(string(data)))
}

func prettyPrint(words []string) {
	for i, word := range words {
		fmt.Printf("%4d %s\n", i, word)
	}
}
