package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Println("Usage: ./hack-assembler <filename.asm>")
		os.Exit(1)
	}

	filename := os.Args[1]
	ext := filepath.Ext(filename)
	outputFile := strings.TrimSuffix(filename, ext) + ".hack"

	parser, err := NewParser(filename)
	if err != nil {
		log.Fatalln(err)
	}

	symbolTable := NewSymbolTable()

	romAddress := 0

	for parser.HasMoreCommands() {
		parser.Advance()

		if parser.CommandType() == LCommand {
			symbolTable.AddEntry(parser.Symbol(), romAddress)
		} else {
			romAddress++
		}
	}

	parser.Reset()

	code := &Code{}
	variableAddress := 16

	var buf bytes.Buffer
	first := true

	writeLine := func(line string) {
		if !first {
			buf.WriteByte('\n')
		}
		buf.WriteString(line)
		first = false
	}

	for parser.HasMoreCommands() {
		parser.Advance()

		switch parser.CommandType() {

		case ACommand:
			symbol := parser.Symbol()

			var address int

			if n, err := strconv.Atoi(symbol); err == nil {
				address = n

			} else if symbolTable.Contains(symbol) {
				address = symbolTable.GetAddress(symbol)

			} else {
				address = variableAddress
				symbolTable.AddEntry(symbol, address)
				variableAddress++
			}

			writeLine(fmt.Sprintf("%016b", address))

		case CCommand:
			comp := code.Comp(parser.Comp())
			dest := code.Dest(parser.Dest())
			jump := code.Jump(parser.Jump())

			writeLine(fmt.Sprintf("111%s%s%s", comp, dest, jump))

		case LCommand:
		}
	}

	if err := os.WriteFile(outputFile, buf.Bytes(), 0644); err != nil {
		log.Fatalln(err)
	}
}
