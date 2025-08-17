package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tyrkinn/stdtest/internal/language/tokenizer"
)

func main() {
	filePath := flag.String("f", ".stdtest", "test file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	tokenizer := tokenizer.New(reader)
	tokens, err := tokenizer.ScanTokens()
	fmt.Println(tokens)
}
