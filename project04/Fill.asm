(LOOP)
    // Check keyboard
    @KBD
    D=M

    // If D != 0, go black
    @BLACK
    D;JNE

    // Otherwise clear screen
    @WHITE
    0;JMP

// --------------------
// Fill screen BLACK
// --------------------
(BLACK)
    @SCREEN
    D=A
    @ADDR
    M=D

(BLACK_LOOP)
    @ADDR
    D=M

    @KBD
    D=D-A
    @LOOP
    D;JEQ

    @ADDR
    A=M
    M=-1

    @ADDR
    M=M+1

    @BLACK_LOOP
    0;JMP

// --------------------
// Fill screen WHITE
// --------------------
(WHITE)
    @SCREEN
    D=A
    @ADDR
    M=D

(WHITE_LOOP)
    @ADDR
    D=M

    @KBD
    D=D-A
    @LOOP
    D;JEQ

    @ADDR
    A=M
    M=0

    @ADDR
    M=M+1

    @WHITE_LOOP
    0;JMP
