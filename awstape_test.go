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
	"io"
	"reflect"
	"testing"
)

var (
	// Scenario 1: valid tape mark block
	tapeMarkBlock = &AWSTapeBlock{
		CurrentLength:  0,
		PreviousLength: 0,
		Flags1:         AWSTapeMark,
		Flags2:         0,
		Data:           []byte{},
	}
	tapeMarkBinary = []byte{0, 0, 0, 0, 0x40, 0}

	// Scenario 2: valid logical record block
	validData      = []byte("testing 1, 2, 3")
	validDataBlock = &AWSTapeBlock{
		CurrentLength:  15,
		PreviousLength: 0,
		Flags1:         AWSRecordStart | AWSRecordEnd,
		Flags2:         0,
		Data:           validData,
	}
	validDataBinary = append([]byte{0xF, 0, 0, 0, 0xA0, 0}, validData...)

	// Scenario 3: invalid block - header too short
	invalidHeaderBinary = []byte{0x0, 0x0, 0x0, 0x0}

	// Scenario 4: invalid block - missing data frame
	missingDataBinary = []byte{0x0A, 0, 0, 0, 0xA0, 0}

	// Scenario 5: invalid block - data frame too short
	lengthMismatchBinary = []byte{0x55, 0xAA, 0, 0, 0xA0, 0, 1, 2, 3, 4}
)

func testMarshal(t *testing.T, testCase *AWSTapeBlock, expected []byte) {
	result := testCase.Marshal()

	if !bytes.Equal(result, expected) {
		t.Errorf("AWSTapeBlock.Marshal failed: result = %+v, expected = %+v", result, expected)
	}
}

func testUnmarshal(t *testing.T, testCase []byte, expected *AWSTapeBlock) {
	reader := bytes.NewReader(testCase)
	result, err := UnmarshalAWSBlock(reader)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("AWSTapeBlock.Unmarshal failed: result = %+v, expected = %+v", result, expected)
	}
}

func testUnmarshalFailure(t *testing.T, testCase []byte, expected error) {
	reader := bytes.NewReader(testCase)

	result, err := UnmarshalAWSBlock(reader)

	if result != nil {
		t.Fatalf("UnmarshalAWSBlock should have returned a nil result (got: %+v)", result)
	}
	if err != expected {
		t.Errorf("UnmarshalAWSBlock should have returned %+v (got: %+v)", expected, err)
	}
}

func TestMarshalTapeMark(t *testing.T) {
	testMarshal(t, tapeMarkBlock, tapeMarkBinary)
}

func TestUnmarshalTapeMark(t *testing.T) {
	testUnmarshal(t, tapeMarkBinary, tapeMarkBlock)
}

func TestMarshalValidData(t *testing.T) {
	testMarshal(t, validDataBlock, validDataBinary)
}

func TestUnmarshalValidData(t *testing.T) {
	testUnmarshal(t, validDataBinary, validDataBlock)
}

func TestUnmarshalShortHeader(t *testing.T) {
	testUnmarshalFailure(t, invalidHeaderBinary, io.ErrUnexpectedEOF)
}

func TestUnmarshalNoData(t *testing.T) {
	testUnmarshalFailure(t, missingDataBinary, io.ErrUnexpectedEOF)
}

func TestUnmarshalShortBlock(t *testing.T) {
	testUnmarshalFailure(t, lengthMismatchBinary, io.ErrUnexpectedEOF)
}
