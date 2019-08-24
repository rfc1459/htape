package htape

import (
	"bytes"
	"encoding/binary"
	"io"
)

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

//-----------------------------------------------------------------------------
// An AWS tape image is a sequence of dynamic-length blocks with a 6-byte
// header. The header is defined as follows (multi-byte integers are
// little-endian):
//   1. length of current block (2 bytes, LE)
//   2. length of previous block (2 bytes, LE)
//   3. Flags1 (1 byte)
//   4. Flags2 (1 byte, currently unused)
//
// Flag1 is a bitmask for the following flags
//   1. 0x80 - Start of new record
//   2. 0x40 - Tape mark
//   3. 0x20 - End of record
//-----------------------------------------------------------------------------

var (
	byteOrder = binary.LittleEndian
)

// AWSFlags is a bitfield representing flags for a single AWS block
type AWSFlags uint8

const (
	// AWSRecordEnd indicates the end of a logical record
	AWSRecordEnd AWSFlags = 0x20 << iota
	// AWSTapeMark indicates this block is a tape mark (no data)
	AWSTapeMark
	// AWSRecordStart indicates the start of a logical record
	AWSRecordStart
)

// AWSTapeBlock is a single block of an AWS tape image
type AWSTapeBlock struct {
	// CurrentLength is the length of the current block
	CurrentLength uint16
	// PreviousLength is the length of the previous block
	PreviousLength uint16

	// Flags1 is the field for standard AWS flags (see AWSFlags)
	Flags1 AWSFlags
	// Flags2 is a reserved field for future flags
	Flags2 uint8

	// Data for current block
	Data []byte
}

// Marshal create a binary representation of this block
func (b *AWSTapeBlock) Marshal() []byte {
	buf := new(bytes.Buffer)
	header := make([]byte, 6)
	byteOrder.PutUint16(header[0:2], b.CurrentLength)
	byteOrder.PutUint16(header[2:4], b.PreviousLength)
	header[4] = byte(b.Flags1)
	header[5] = b.Flags2
	buf.Write(header)
	buf.Write(b.Data)
	rv := make([]byte, buf.Len())
	copy(rv, buf.Bytes())
	return rv
}

// UnmarshalAWSBlock reads an AWS tape block from a stream
func UnmarshalAWSBlock(in io.Reader) (*AWSTapeBlock, error) {
	// Read the block header
	header := make([]byte, 6)
	n, err := in.Read(header)
	if err != nil && err != io.EOF {
		return nil, err
	}
	if n < len(header) {
		return nil, io.ErrUnexpectedEOF
	}

	// Deserialize the block header
	block := &AWSTapeBlock{
		CurrentLength:  byteOrder.Uint16(header[0:2]),
		PreviousLength: byteOrder.Uint16(header[2:4]),
		Flags1:         AWSFlags(header[4]),
		Flags2:         header[5],
		Data:           []byte{},
	}

	// TODO: perform header sanity checks

	// Read data (if any)
	if block.CurrentLength > 0 {
		dlen := int(block.CurrentLength)
		data := make([]byte, dlen)
		n, err = in.Read(data)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n < dlen {
			return nil, io.ErrUnexpectedEOF
		}
		block.Data = data
	}
	return block, nil
}
