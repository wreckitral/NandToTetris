package main

import (
	"fmt"
	"os"
	"strconv"
)

type CodeWriter struct {
	file *os.File
	fileName string
	labelCounter int
}

func NewCodeWriter(filename string) (*CodeWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	return &CodeWriter{
		file: file,
		fileName: filename,
		labelCounter: 0,
	}, nil
}

func (cw *CodeWriter) Close() error {
	return cw.file.Close()
}

func (cw *CodeWriter) emit(lines ...string) {
	for _, line := range lines {
		fmt.Fprintln(cw.file, line)
	}
}

var binaryOps = map[string]string{
	"add": "M=D+M",
	"sub": "M=M-D",
	"and": "M=D&M",
	"or":  "M=D|M",
}

var unaryOps = map[string]string{
	"neg": "M=-M",
	"not": "M=!M",
}

var jumpOps = map[string]string{
	"eq": "D;JEQ",
	"gt": "D;JGT",
	"lt": "D;JLT",
}


func (cw *CodeWriter) WriteArithmetic(command string) error {
	if op, ok := binaryOps[command]; ok {
		cw.emit(
			"@SP",
			"AM=M-1",
			"D=M",
			"@SP",
			"A=M-1",
			op,
			"",
		)
		return nil
	}

	if op, ok := unaryOps[command]; ok {
		cw.emit(
			"@SP",
			"A=M-1",
			op,
			"",
		)
		return nil
	}

	if jump, ok := jumpOps[command]; ok {
		jumpLabel := "CMP_TRUE" + strconv.Itoa(cw.labelCounter)
		cw.labelCounter++

		cw.emit(
			"@SP",
			"AM=M-1",
			"D=M",
			"@SP",
			"A=M-1",
			"D=M-D",
			"M=-1",
			"@"+jumpLabel,
			jump,
			"@SP",
			"A=M-1",
			"M=0",
			"("+jumpLabel+")",
			"",
		)
		return nil
	}

	return fmt.Errorf("unknown arithmetic command %q", command)
}

var segmentMap = map[string]string{
	"local":    "@LCL",
	"argument": "@ARG",
	"this":     "@THIS",
	"that":     "@THAT",
	"pointer":  "@3",
	"temp":     "@5",
}

func (cw *CodeWriter) WritePushPop(command CommandType, segment string, index int) error {
	idxStr := strconv.Itoa(index)

	switch command {
	case CPush:
		switch segment {
		case "constant":
			cw.emit(
				"@"+idxStr,
				"D=A",
				"@SP",
				"AM=M+1",
				"A=A-1",
				"M=D",
			)
		case "local", "argument", "this", "that", "temp", "pointer":
			cw.emit(
				"@"+idxStr,
				"D=A",
			)

			if segment == "temp" || segment == "pointer" {
				cw.emit(segmentMap[segment])
			} else {
				cw.emit(
					segmentMap[segment],
					"A=M",
				)
			}

			cw.emit(
				"A=D+A",
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			)
		case "static":
			cw.emit(
				fmt.Sprintf("@%s.%d", cw.fileName, index),
				"D=M",
				"@SP",
				"A=M",
				"M=D",
				"@SP",
				"M=M+1",
			)
		default:
			return fmt.Errorf("unexpected push segment: %q", segment)
		}

	case CPop:
		switch segment {
		case "constant":
			return fmt.Errorf("cannot pop constant segment")
		case "local", "argument", "this", "that", "temp", "pointer":
			cw.emit(
				"@"+idxStr,
				"D=A",
			)

			if segment == "temp" || segment == "pointer" {
				cw.emit(segmentMap[segment])
			} else {
				cw.emit(
					segmentMap[segment],
					"A=M",
				)
			}

			cw.emit(
				"D=D+A",
				"@R13",
				"M=D",
				"@SP",
				"AM=M-1",
				"D=M",
				"@R13",
				"A=M",
				"M=D",
			)
		case "static":
			cw.emit(
				"@SP",
				"AM=M-1",
				"D=M",
				fmt.Sprintf("@%s.%d", cw.fileName, index),
				"M=D",
			)
		default:
			return fmt.Errorf("unexpected pop segment: %q", segment)
		}
	default:
		return fmt.Errorf("unexpected command type: %q", command)
	}

	cw.emit("")
	return nil
}
