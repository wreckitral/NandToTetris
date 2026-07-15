// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Mult.asm

// Computes R2 = R0 * R1

// Initialize result
@R2
M=0

// If R0 == 0, done
@R0
D=M
@END
D;JEQ

// If R1 == 0, done
@R1
D=M
@END
D;JEQ

// COUNT = R1
@R1
D=M
@COUNT
M=D

(LOOP)
    // if COUNT == 0, end
    @COUNT
    D=M
    @END
    D;JEQ

    // R2 = R2 + R0
    @R0
    D=M
    @R2
    M=D+M

    // COUNT--
    @COUNT
    M=M-1

    // Repeat
    @LOOP
    0;JMP

(END)
    @END
    0;JMP
