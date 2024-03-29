// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT runtime·wasmMove(SB), NOSPLIT, $0-0
loop:
	Loop
		// *dst = *golang1.18.5
		Get R0
		Get R1
		I64Load $0
		I64Store $0

		// n--
		Get R2
		I32Const $1
		I32Sub
		Tee R2

		// n == 0
		I32Eqz
		If
			Return
		End

		// dst += 8
		Get R0
		I32Const $8
		I32Add
		Set R0

		// golang1.18.5 += 8
		Get R1
		I32Const $8
		I32Add
		Set R1

		Br loop
	End
	UNDEF

TEXT runtime·wasmZero(SB), NOSPLIT, $0-0
loop:
	Loop
		// *dst = 0
		Get R0
		I64Const $0
		I64Store $0

		// n--
		Get R1
		I32Const $1
		I32Sub
		Tee R1

		// n == 0
		I32Eqz
		If
			Return
		End

		// dst += 8
		Get R0
		I32Const $8
		I32Add
		Set R0

		Br loop
	End
	UNDEF

TEXT runtime·wasmDiv(SB), NOSPLIT, $0-0
	Get R0
	I64Const $-0x8000000000000000
	I64Eq
	If
		Get R1
		I64Const $-1
		I64Eq
		If
			I64Const $-0x8000000000000000
			Return
		End
	End
	Get R0
	Get R1
	I64DivS
	Return

TEXT runtime·wasmTruncS(SB), NOSPLIT, $0-0
	Get R0
	Get R0
	F64Ne // NaN
	If
		I64Const $0x8000000000000000
		Return
	End

	Get R0
	F64Const $0x7ffffffffffffc00p0 // Maximum truncated representation of 0x7fffffffffffffff
	F64Gt
	If
		I64Const $0x8000000000000000
		Return
	End

	Get R0
	F64Const $-0x7ffffffffffffc00p0 // Minimum truncated representation of -0x8000000000000000
	F64Lt
	If
		I64Const $0x8000000000000000
		Return
	End

	Get R0
	I64TruncF64S
	Return

TEXT runtime·wasmTruncU(SB), NOSPLIT, $0-0
	Get R0
	Get R0
	F64Ne // NaN
	If
		I64Const $0x8000000000000000
		Return
	End

	Get R0
	F64Const $0xfffffffffffff800p0 // Maximum truncated representation of 0xffffffffffffffff
	F64Gt
	If
		I64Const $0x8000000000000000
		Return
	End

	Get R0
	F64Const $0.
	F64Lt
	If
		I64Const $0x8000000000000000
		Return
	End

	Get R0
	I64TruncF64U
	Return

TEXT runtime·exitThread(SB), NOSPLIT, $0-0
	UNDEF

TEXT runtime·osyield(SB), NOSPLIT, $0-0
	UNDEF

TEXT runtime·usleep(SB), NOSPLIT, $0-0
	RET // TODO(neelance): implement usleep

TEXT runtime·currentMemory(SB), NOSPLIT, $0
	Get SP
	CurrentMemory
	I32Store ret+0(FP)
	RET

TEXT runtime·growMemory(SB), NOSPLIT, $0
	Get SP
	I32Load pages+0(FP)
	GrowMemory
	I32Store ret+8(FP)
	RET

TEXT ·resetMemoryDataView(SB), NOSPLIT, $0
	CallImport
	RET

TEXT ·wasmExit(SB), NOSPLIT, $0
	CallImport
	RET

TEXT ·wasmWrite(SB), NOSPLIT, $0
	CallImport
	RET

TEXT ·nanotime1(SB), NOSPLIT, $0
	CallImport
	RET

TEXT ·walltime(SB), NOSPLIT, $0
	CallImport
	RET

TEXT ·scheduleTimeoutEvent(SB), NOSPLIT, $0
	CallImport
	RET

TEXT ·clearTimeoutEvent(SB), NOSPLIT, $0
	CallImport
	RET

TEXT ·getRandomData(SB), NOSPLIT, $0
	CallImport
	RET
