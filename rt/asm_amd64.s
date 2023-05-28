#include "textflag.h"

TEXT ·getgptr(SB),NOSPLIT,$0-8
    MOVQ (TLS), AX
    MOVQ AX, ret+0(FP)
    RET
