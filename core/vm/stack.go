// Copyright 2014 The go-ethereum Authors
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
	"sync"

	"github.com/holiman/uint256"
)

var stackPool = sync.Pool{
	New: func() interface{} {
		return &Stack{data: make([]uint256.Int, 0, 16)}
	},
}

// Stack is an object for basic stack operations. Items popped to the stack are
// expected to be changed and modified. stack does not take care of adding newly
// initialized objects.
type Stack struct {
	data []uint256.Int
}

func newstack() *Stack {
	return stackPool.Get().(*Stack)
}

func returnStack(s *Stack) {
	s.data = s.data[:0]
	stackPool.Put(s)
}

// Data returns the underlying uint256.Int array.
func (st *Stack) Data() []uint256.Int {
	return st.data
}

func (st *Stack) push(d *uint256.Int) {
	// NOTE push limit (1024) is checked in baseCheck
	st.data = append(st.data, *d)
}

func (st *Stack) pop() (ret uint256.Int) {
	ret = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}

func (st *Stack) pop2() (pop0, pop1 uint256.Int) {
	pop1 = st.data[len(st.data)-2]
	pop0 = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-2]
	return
}

func (st *Stack) pop3() (pop0, pop1, pop2 uint256.Int) {
	pop2 = st.data[len(st.data)-3]
	pop1 = st.data[len(st.data)-2]
	pop0 = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-3]
	return
}

func (st *Stack) pop4() (pop0, pop1, pop2, pop3 uint256.Int) {
	pop3 = st.data[len(st.data)-4]
	pop2 = st.data[len(st.data)-3]
	pop1 = st.data[len(st.data)-2]
	pop0 = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-4]
	return
}

func (st *Stack) popPeek() (pop uint256.Int, peek *uint256.Int) {
	peek = &st.data[len(st.data)-2]
	pop = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}

func (st *Stack) pop2Peek() (pop0, pop1 uint256.Int, peek *uint256.Int) {
	peek = &st.data[len(st.data)-3]
	pop1 = st.data[len(st.data)-2]
	pop0 = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-2]
	return
}

func (st *Stack) pop3Peek() (pop0, pop1, pop2 uint256.Int, peek *uint256.Int) {
	peek = &st.data[len(st.data)-4]
	pop2 = st.data[len(st.data)-3]
	pop1 = st.data[len(st.data)-2]
	pop0 = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-3]
	return
}

func (st *Stack) pop5Peek() (pop0, pop1, pop2, pop3, pop4 uint256.Int, peek *uint256.Int) {
	peek = &st.data[len(st.data)-6]
	pop4 = st.data[len(st.data)-5]
	pop3 = st.data[len(st.data)-4]
	pop2 = st.data[len(st.data)-3]
	pop1 = st.data[len(st.data)-2]
	pop0 = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-5]
	return
}

func (st *Stack) pop6Peek() (pop0, pop1, pop2, pop3, pop4, pop5 uint256.Int, peek *uint256.Int) {
	peek = &st.data[len(st.data)-7]
	pop5 = st.data[len(st.data)-6]
	pop4 = st.data[len(st.data)-5]
	pop3 = st.data[len(st.data)-4]
	pop2 = st.data[len(st.data)-3]
	pop1 = st.data[len(st.data)-2]
	pop0 = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-6]
	return
}

func (st *Stack) len() int {
	return len(st.data)
}

func (st *Stack) swap1()  { tmp := st.data[st.len()-2:]; tmp[0], tmp[1] = tmp[1], tmp[0] }
func (st *Stack) swap2()  { tmp := st.data[st.len()-3:]; tmp[0], tmp[2] = tmp[2], tmp[0] }
func (st *Stack) swap3()  { tmp := st.data[st.len()-4:]; tmp[0], tmp[3] = tmp[3], tmp[0] }
func (st *Stack) swap4()  { tmp := st.data[st.len()-5:]; tmp[0], tmp[4] = tmp[4], tmp[0] }
func (st *Stack) swap5()  { tmp := st.data[st.len()-6:]; tmp[0], tmp[5] = tmp[5], tmp[0] }
func (st *Stack) swap6()  { tmp := st.data[st.len()-7:]; tmp[0], tmp[6] = tmp[6], tmp[0] }
func (st *Stack) swap7()  { tmp := st.data[st.len()-8:]; tmp[0], tmp[7] = tmp[7], tmp[0] }
func (st *Stack) swap8()  { tmp := st.data[st.len()-9:]; tmp[0], tmp[8] = tmp[8], tmp[0] }
func (st *Stack) swap9()  { tmp := st.data[st.len()-10:]; tmp[0], tmp[9] = tmp[9], tmp[0] }
func (st *Stack) swap10() { tmp := st.data[st.len()-11:]; tmp[0], tmp[10] = tmp[10], tmp[0] }
func (st *Stack) swap11() { tmp := st.data[st.len()-12:]; tmp[0], tmp[11] = tmp[11], tmp[0] }
func (st *Stack) swap12() { tmp := st.data[st.len()-13:]; tmp[0], tmp[12] = tmp[12], tmp[0] }
func (st *Stack) swap13() { tmp := st.data[st.len()-14:]; tmp[0], tmp[13] = tmp[13], tmp[0] }
func (st *Stack) swap14() { tmp := st.data[st.len()-15:]; tmp[0], tmp[14] = tmp[14], tmp[0] }
func (st *Stack) swap15() { tmp := st.data[st.len()-16:]; tmp[0], tmp[15] = tmp[15], tmp[0] }
func (st *Stack) swap16() { tmp := st.data[st.len()-17:]; tmp[0], tmp[16] = tmp[16], tmp[0] }

func (st *Stack) dup(n int) {
	st.push(&st.data[st.len()-n])
}

func (st *Stack) peek() *uint256.Int {
	return &st.data[st.len()-1]
}

// Back returns the n'th item in stack
func (st *Stack) Back(n int) *uint256.Int {
	return &st.data[st.len()-n-1]
}
