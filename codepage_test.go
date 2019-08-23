package htape

// Copyright (c) 2019, Matteo Panella
// All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:

// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.

// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.

// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

import (
	"bytes"
	"testing"

	"golang.org/x/text/transform"
)

var (
	asciiBytes  = []byte("htape Transform test")
	ebcdicBytes = []byte{0x88, 0xa3, 0x81, 0x97, 0x85, 0x40, 0xe3, 0x99, 0x81, 0x95, 0xa2, 0x86, 0x96, 0x99, 0x94, 0x40, 0xa3, 0x85, 0xa2, 0xa3}
)

func TestASCIIToEBCDIC(t *testing.T) {
	res, n, err := transform.Bytes(asciiToEBCDIC, asciiBytes)
	if err != nil {
		t.Error(err)
	}

	if n != len(asciiBytes) || !bytes.Equal(res, ebcdicBytes) {
		t.Errorf("asciiToEBCDIC.Transform failed: result = %+v, expected = %+v", res, ebcdicBytes)
	}
}

func TestEBCDICToASCII(t *testing.T) {
	res, n, err := transform.Bytes(ebcdicToASCII, ebcdicBytes)
	if err != nil {
		t.Error(err)
	}

	if n != len(ebcdicBytes) || !bytes.Equal(res, asciiBytes) {
		t.Errorf("asciiToEBCDIC.Transform failed: result = %+v, expected = %+v", res, ebcdicBytes)
	}
}

func TestDefaultMapsIdentity(t *testing.T) {
	for i := 0; i < 256; i++ {
		expected := byte(i)
		eb := asciiToEBCDICMap[expected]
		res := ebcdicToASCIIMap[eb]
		// Special case for ASCII 0xb4 and EBCDIC 0x15: the backward conversion is *not* symmetric
		if expected == 0xb4 && eb == 0x15 {
			expected = 0x0a
		}
		if res != expected {
			t.Errorf("Identity mapping failed: result = %#02x (EBCDIC: %#02x), expected = %#02x", res, eb, expected)
		}
	}
}
