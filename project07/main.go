package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 2 || filepath.Ext(os.Args[1]) != ".vm" {
		fmt.Println("Usage: go run . [file.vm]")
		os.Exit(1)
	}

	inFile := os.Args[1]

	outFile := strings.TrimSuffix(inFile, ".vm") + ".asm"

	baseName := strings.TrimSuffix(filepath.Base(inFile), ".vm")

	parser, err := NewParser(inFile)
	if err != nil {
		log.Fatalf("Failed to create parser: %v", err)
	}

	cw, err := NewCodeWriter(outFile)
	if err != nil {
		log.Fatalf("Failed to create code writer: %v", err)
	}
	defer cw.Close()

	cw.fileName = baseName

	for parser.HasMoreCommands() {
		parser.Advance()

		cw.emit("// " + parser.Current())

		cmdType := parser.CommandType()
		switch cmdType {
		case CArithmetic:
			if err := cw.WriteArithmetic(parser.Arg1()); err != nil {
				log.Fatalf("Error writing arithmetic: %v", err)
			}
		case CPush, CPop:
			if err := cw.WritePushPop(cmdType, parser.Arg1(), parser.Arg2()); err != nil {
				log.Fatalf("Error writing push/pop: %v", err)
			}
		case CUnknown:
			log.Fatalf("Unknown command: %s", parser.Current())
		}
	}
}
