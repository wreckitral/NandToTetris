package main

import (
	"os"
	"strings"
)

type Parser struct {
	commands []string
	current  int
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

func (p *Parser) Reset() {
	p.current = -1
}

type CommandType string

const (
	ACommand CommandType = "A"
	CCommand CommandType = "C"
	LCommand CommandType = "L"
)

func (p *Parser) CommandType() CommandType {
	current := p.Current()

	if strings.HasPrefix(current, "@") {
		return ACommand
	}

	if strings.HasPrefix(current, "(") {
		return LCommand
	}

	return CCommand
}

func (p *Parser) Symbol() string {
	line := strings.TrimPrefix(p.Current(), "@")
	line = strings.TrimPrefix(line, "(")
	line = strings.TrimSuffix(line, ")")

	return line
}

func (p *Parser) Dest() string {
	line, _, found := strings.Cut(p.Current(), "=")
	if !found {
		return ""
	}

	return line
}

func (p *Parser) Comp() string {
	line, _, _ := strings.Cut(p.Current(), ";")
	dest, comp, found := strings.Cut(line, "=")
	if !found {
		return dest
	}

	return comp
}

func (p *Parser) Jump() string {
	_, line , _ := strings.Cut(p.Current(), ";")

	return line
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
