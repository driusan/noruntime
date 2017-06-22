TEXT main+0(SB), $0
    MOVL    $33, DI      // arg 1 exit status
    MOVL    $(0x2000000+1), AX  // syscall entry
    SYSCALL
    MOVL    $0xf1, 0xf1  // crash
    RET

