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
	"golang.org/x/text/transform"
)

type asciiToEBCDICTransformer struct {
	transform.NopResetter
}

type ebcdicToASCIITransformer struct {
	transform.NopResetter
}

var (
	asciiToEBCDIC transform.Transformer = asciiToEBCDICTransformer{}
	ebcdicToASCII transform.Transformer = ebcdicToASCIITransformer{}
)

// Transform translates a stream of bytes from ASCII to EBCDIC
func (asciiToEBCDICTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	nDst, nSrc, err = doTransform(dst, src, asciiToEBCDICMap, atEOF)
	return
}

// Transform translates a stream of bytes from EBCDIC to ASCII
func (ebcdicToASCIITransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	nDst, nSrc, err = doTransform(dst, src, ebcdicToASCIIMap, atEOF)
	return
}

// TODO: generalize transformer for other codepages - see hercules/codepage.c
// for codepages we should (probably) support

func doTransform(dst, src, mapping []byte, atEOF bool) (nDst, nSrc int, err error) {
	for _, c := range src {
		if nDst == len(dst) {
			err = transform.ErrShortDst
			break
		}
		dst[nDst] = mapping[c]
		nDst++
		nSrc++
	}
	return
}
