// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/exp/slices"
)

func TestBitvec(t *testing.T) {
	tests := []struct {
		Code []byte
		Want bitvec
	}{
		{Code: []byte{}, Want: bitvec{0, 0}},
		{Code: []byte{byte(PUSH1), 0x01, 0x01, 0x01}, Want: bitvec{0b00000000_00000000_00000000_00000010, 0}},
		{Code: []byte{byte(PUSH2), 0x01, 0x01, 0x01}, Want: bitvec{0b00000000_00000000_00000000_00000110, 0}},
		{
			Code: []byte{byte(PUSH1), byte(PUSH1), byte(PUSH1), byte(PUSH1)},
			Want: bitvec{0b00000000_00000000_00000000_00001010, 0},
		},
		{
			Code: []byte{0x00, byte(PUSH1), 0x00, byte(PUSH1), 0x00, byte(PUSH1), 0x00, byte(PUSH1)},
			Want: bitvec{0b00000000_00000000_00000001_01010100, 0},
		},
		{
			Code: []byte{byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), byte(PUSH8), 0x01, 0x01, 0x01},
			Want: bitvec{0b00000000_00000000_00000001_11111110, 0},
		},
		{
			Code: []byte{0x01, 0x01, 0x01, 0x01, 0x01, byte(PUSH2), 0x01, 0x01, 0x01, 0x01, 0x01},
			Want: bitvec{0b00000000_00000000_00000000_11000000, 0},
		},

		{
			Code: []byte{byte(PUSH3), 0x01, 0x01, 0x01, byte(PUSH1), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			Want: bitvec{0b00000000_00000000_00000000_00101110, 0},
		},
		{
			Code: []byte{0x01, byte(PUSH8), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			Want: bitvec{0b00000000_00000000_00000011_11111100, 0},
		},

		{
			Code: []byte{byte(PUSH16), 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01},
			Want: bitvec{0b00000000_00000001_11111111_11111110, 0},
		},
		{
			Code: []byte{byte(PUSH8), 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, byte(PUSH1), 0x01},
			Want: bitvec{0b_00000000_00000000_00000101_11111110, 0},
		},
		{Code: []byte{byte(PUSH32)}, Want: bitvec{0b11111111_11111111_11111111_11111110, 0b00000000_00000000_00000000_00000001}},
		{
			Code: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, byte(PUSH2), 0xff,
				0xff,
			},
			Want: bitvec{0b10000000_00000000_00000000_00000000, 0b00000000_00000000_00000000_00000001, 0},
		},
		{
			Code: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, byte(PUSH32),
			},
			Want: bitvec{0b00000000_00000000_00000000_00000000, 0b11111111_11111111_11111111_11111111, 0},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := newCodeBitVec(test.Code)
			if !slices.Equal(test.Want, got) {
				t.Fatalf("(-want +got)\n- %32b\n+ %32b\n", test.Want, got)
			}
		})
	}
}

const analysisCodeSize = 1200 * 1024

func BenchmarkJumpdestAnalysis_1200k(bench *testing.B) {
	// 1.4 ms
	code := make([]byte, analysisCodeSize)
	bench.SetBytes(analysisCodeSize)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		newCodeBitVec(code)
	}
}

func BenchmarkJumpdestAnalysis_rand(b *testing.B) {
	code := make([]byte, analysisCodeSize)
	rand.Read(code)

	bv := newCodeBitVec(code)
	b.SetBytes(int64(len(code)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bv.codeBitvecInternal(code)
	}
}

var (
	//go:embed testdata/weth9.bytecode
	hexCodeWETH9 string
	codeWETH9    = common.FromHex(hexCodeWETH9)
)

func BenchmarkJumpdestAnalysis_weth9(b *testing.B) {
	bv := newCodeBitVec(codeWETH9)
	b.SetBytes(int64(len(codeWETH9)))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bv.codeBitvecInternal(codeWETH9)
	}
}

func BenchmarkJumpdestHashing_1200k(bench *testing.B) {
	// 4 ms
	code := make([]byte, analysisCodeSize)
	bench.SetBytes(analysisCodeSize)
	bench.ResetTimer()
	for i := 0; i < bench.N; i++ {
		crypto.Keccak256Hash(code)
	}
}

func BenchmarkJumpdestOpAnalysis(b *testing.B) {
	var op OpCode
	bencher := func(b *testing.B) {
		code := bytes.Repeat([]byte{byte(op)}, analysisCodeSize)
		bv := newCodeBitVec(code)
		b.SetBytes(analysisCodeSize)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := range bv {
				bv[j] = 0
			}
			bv.codeBitvecInternal(code)
		}
	}
	for op = PUSH1; op <= PUSH32; op++ {
		b.Run(op.String(), bencher)
	}
	op = JUMPDEST
	b.Run(op.String(), bencher)
	op = STOP
	b.Run(op.String(), bencher)
}
