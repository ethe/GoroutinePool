package gp

import (
	"errors"
	"sync/atomic"
	"unsafe"
)

var (
	EmptyQueue = errors.New("queue emtpy")
)

type node struct {
	inner interface{}
	next  *node
}

type Queue struct {
	head *node
	tail *node
}

func NewQueue() *Queue {
	dummy := &node{nil, nil}
	return &Queue{
		dummy,
		dummy,
	}
}

func (q *Queue) Put(v interface{}) {
	n := &node{v, nil}
	for {
		tail := q.tail
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail.next)), nil, unsafe.Pointer(n)) {
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), unsafe.Pointer(tail), (unsafe.Pointer)(q.tail.next))
			return
		}
	}
}

func (q *Queue) Get() (interface{}, error) {
	for {
		p := q.head
		if p.next == nil {
			return nil, EmptyQueue
		}
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), unsafe.Pointer(p), (unsafe.Pointer)(q.head.next)) {
			return p.next.inner, nil
		}
	}
}
