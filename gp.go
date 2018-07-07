// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

// NOTICE File was modified by Tzuhsing Gwo

package gp

import (
	"sync/atomic"
	"time"
)

// Pool is a struct to represent goroutine pool.
type Pool struct {
	*Queue
	idleTimeout time.Duration
}

// goroutine is actually a background goroutine, with a channel binded for communication.
type goroutine struct {
	ch     chan func()
	status int32
}

const (
	statusIdle  int32 = 0
	statusInUse int32 = 1
	statusDead  int32 = 2
)

// New returns a new *Pool object.
func New(idleTimeout time.Duration) *Pool {
	pool := &Pool{
		NewQueue(),
		idleTimeout,
	}
	return pool
}

// Go works like go func(), but goroutines are pooled for reusing.
// This strategy can avoid runtime.morestack, because pooled goroutine is already enlarged.
func (pool *Pool) Go(f func()) {
	for {
		g := pool.get()
		if atomic.CompareAndSwapInt32(&g.status, statusIdle, statusInUse) {
			g.ch <- f
			return
		}
		// Status already changed from statusIdle => statusDead, drop it, find next one.
	}
}

func (pool *Pool) get() *goroutine {
	ret, err := pool.Queue.Get()
	if err == EmptyQueue {
		return pool.alloc()
	}
	return ret.(*goroutine)
}

func (pool *Pool) alloc() *goroutine {
	g := &goroutine{
		ch: make(chan func()),
	}
	go g.workLoop(pool)
	return g
}

func (g *goroutine) putInto(pool *Pool) {
	g.status = statusIdle
	pool.Put(g)
}

func (g *goroutine) workLoop(pool *Pool) {
	timer := time.NewTimer(pool.idleTimeout)
	for {
		select {
		case <-timer.C:
			// Check to avoid a corner case that the goroutine is take out from pool,
			// and Get this signal at the same time.
			succ := atomic.CompareAndSwapInt32(&g.status, statusIdle, statusDead)
			if succ {
				return
			}
		case work := <-g.ch:
			work()
			// Put g back to the pool.
			// This is the normal usage for a resource pool:
			//
			//     obj := pool.Get()
			//     use(obj)
			//     pool.putInto(obj)
			//
			// But when goroutine is used as a resource, we can't pool.putInto() immediately,
			// because the resource(goroutine) maybe still in use.
			// So, putInto back resource is done here,  when the goroutine finish its work.
			g.putInto(pool)
		}
		timer.Reset(pool.idleTimeout)
	}
}