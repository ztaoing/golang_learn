
    TEXT _rt0_amd64(SB),NOSPLIT,$-8
    MOVQ 0(SP), DI // argc
    LEAQ 8(SP), SI // argv
    JMP runtime·rt0_go(SB)

该方法将程序输入的 argc 和 argv 从内存移动到寄存器中。
栈指针（SP）的前两个值分别是 argc 和 argv，其对应参数的数量和具体各参数的值。

参数的数量+参数的值：这种数据的排列方式是比较常见的