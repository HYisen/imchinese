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

func prettyPrint(candidates []finder.Candidate) {
	for i, candidate := range candidates {
		fmt.Printf("%4d %s 「%s」 %s\n", i, candidate.Word, candidate.Line, candidate.Path)
	}
}
