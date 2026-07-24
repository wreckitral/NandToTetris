package main

import (
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	commands []string
	current int
}

func NewParser(filename string) (*Parser, error) {
	parser := &Parser{
		current: -1,
	}
	if err := parser.readFile(filename); err != nil {
		return nil, err
	}

	return parser, nil
}

func (p *Parser) HasMoreCommands() bool {
	return p.current+1 < len(p.commands)
}

func (p *Parser) Advance() {
	p.current++
}

func (p *Parser) Current() string {
	return p.commands[p.current]
}

type CommandType string

const (
	CUnknown CommandType = "UNKNOWN"
	CArithmetic CommandType = "C_ARITHMETIC"
	CPush CommandType = "C_PUSH"
	CPop CommandType = "C_POP"
)

func (p *Parser) CommandType() CommandType {
	cmd := strings.Fields(p.Current())[0]

	switch cmd {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		return CArithmetic

	case "push":
		return CPush

	case "pop":
		return CPop

	default:
		return CUnknown
	}
}

func (p *Parser) Arg1() string {
	fields := strings.Fields(p.Current())

	if p.CommandType() == CArithmetic {
		return fields[0]
	}

	return fields[1]
}

func (p *Parser) Arg2() int {
	fields := strings.Fields(p.Current())

	n, _ := strconv.Atoi(fields[2])

	return n
}

func (p *Parser) readFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	p.commands = nil

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		cleanedLine := cleanLine(line)

		if cleanedLine == "" {
			continue
		}

		p.commands = append(p.commands, cleanedLine)
	}

	return nil
}

func cleanLine(line string) string {
	line, _, _ = strings.Cut(line, "//")
	line = strings.TrimSpace(line)

	return line
}
