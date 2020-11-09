package GoroutinePool

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
		tail := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)))
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&(*node)(tail).next)), nil, unsafe.Pointer(n)) {
			tailNext := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&(*node)(tail).next)))
			atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.tail)), tail, tailNext)
			return
		}
	}
}

func (q *Queue) Get() (interface{}, error) {
	for {
		head := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)))
		headNext := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&(*node)(head).next)))
		if (*node)(headNext) == nil {
			return nil, EmptyQueue
		}
		if atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&q.head)), head, headNext) {
			return (*node)(headNext).inner, nil
		}
	}
}
