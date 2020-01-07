/**
 * Created by shangyindong
 * Date 2020/1/2 2:09 下午
 **/
package queue

import (
	"runtime"
	"sync/atomic"
)

func roundUp(v uint64) uint64 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v |= v >> 32
	v++
	return v
}

type node struct {
	position uint64
	data     interface{}
}

type nodes []*node

type RingBuffer struct {
	_padding0      [8]uint64
	cursor        uint64
	_padding1      [8]uint64
	tail          uint64
	_padding2      [8]uint64
	mask    uint64
	_padding3      [8]uint64
	nodes          nodes
}

func (rb *RingBuffer) init(size uint64) {
	size = roundUp(size)
	rb.nodes = make(nodes, size)
	for i := uint64(0); i < size; i++ {
		rb.nodes[i] = &node{position: i}
	}
	rb.mask = size - 1
}


func (rb *RingBuffer) Put(item interface{}) (bool, error) {
	var n *node
	pos := atomic.LoadUint64(&rb.tail)
L:
	for {
		n = rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := seq - pos; {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.tail, pos, pos+1) {
				break L
			}
		case dif < 0:
			panic(`occur err when put.`)
		default:
			pos = atomic.LoadUint64(&rb.tail)
		}

		runtime.Gosched()
	}

	n.data = item
	atomic.StoreUint64(&n.position, pos+1)
	return true, nil
}



func (rb *RingBuffer) Get() (interface{}, error) {
	var (
		n     *node
		pos   = atomic.LoadUint64(&rb.cursor)
	)
L:
	for {

		n = rb.nodes[pos&rb.mask]
		seq := atomic.LoadUint64(&n.position)
		switch dif := seq - (pos + 1); {
		case dif == 0:
			if atomic.CompareAndSwapUint64(&rb.cursor, pos, pos+1) {
				break L
			}
		case dif < 0:
			panic(`occur err when get.`)
		default:
			pos = atomic.LoadUint64(&rb.cursor)
		}
		runtime.Gosched()
	}
	data := n.data
	n.data = nil
	atomic.StoreUint64(&n.position, pos+rb.mask+1)
	return data, nil
}



func NewRingBuffer(size uint64) *RingBuffer {
	rb := &RingBuffer{}
	rb.init(size)
	return rb
}
