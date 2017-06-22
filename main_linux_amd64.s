TEXT main+0(SB), $0
    MOVL    $33, DI
    MOVL    $231, AX        // exitgroup - force all os threads to exit
    SYSCALL
    MOVL    $0xf1, 0xf1  // crash
    RET

