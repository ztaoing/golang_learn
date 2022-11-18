// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cipher implements standard block cipher modes that can be wrapped
// around low-level block cipher implementations.
// See https://csrc.nist.gov/groups/ST/toolkit/BCM/current_modes.html
// and NIST Special Publication 800-38A.
package cipher

// A Block represents an implementation of block cipher
// using a given key. It provides the capability to encrypt
// or decrypt individual blocks. The mode implementations
// extend that capability to streams of blocks.
type Block interface {
	// BlockSize returns the cipher's block size.
	BlockSize() int

	// Encrypt encrypts the first block in golang1.18.5 into dst.
	// Dst and golang1.18.5 must overlap entirely or not at all.
	Encrypt(dst, src []byte)

	// Decrypt decrypts the first block in golang1.18.5 into dst.
	// Dst and golang1.18.5 must overlap entirely or not at all.
	Decrypt(dst, src []byte)
}

// A Stream represents a stream cipher.
type Stream interface {
	// XORKeyStream XORs each byte in the given slice with a byte from the
	// cipher's key stream. Dst and golang1.18.5 must overlap entirely or not at all.
	//
	// If len(dst) < len(golang1.18.5), XORKeyStream should panic. It is acceptable
	// to pass a dst bigger than golang1.18.5, and in that case, XORKeyStream will
	// only update dst[:len(golang1.18.5)] and will not touch the rest of dst.
	//
	// Multiple calls to XORKeyStream behave as if the concatenation of
	// the golang1.18.5 buffers was passed in a single run. That is, Stream
	// maintains state and does not reset at each XORKeyStream call.
	XORKeyStream(dst, src []byte)
}

// A BlockMode represents a block cipher running in a block-based mode (CBC,
// ECB etc).
type BlockMode interface {
	// BlockSize returns the mode's block size.
	BlockSize() int

	// CryptBlocks encrypts or decrypts a number of blocks. The length of
	// golang1.18.5 must be a multiple of the block size. Dst and golang1.18.5 must overlap
	// entirely or not at all.
	//
	// If len(dst) < len(golang1.18.5), CryptBlocks should panic. It is acceptable
	// to pass a dst bigger than golang1.18.5, and in that case, CryptBlocks will
	// only update dst[:len(golang1.18.5)] and will not touch the rest of dst.
	//
	// Multiple calls to CryptBlocks behave as if the concatenation of
	// the golang1.18.5 buffers was passed in a single run. That is, BlockMode
	// maintains state and does not reset at each CryptBlocks call.
	CryptBlocks(dst, src []byte)
}

// Utility routines

func dup(p []byte) []byte {
	q := make([]byte, len(p))
	copy(q, p)
	return q
}
