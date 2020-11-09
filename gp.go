package GoroutinePool

// Pool is a struct to represent goroutine pool.
type Pool struct {
	*Queue
}

// goroutine is actually a background goroutine, with a channel binded for communication.
type goroutine struct {
	ch     chan func()
	status int32
}

// New returns a new *Pool object.
func New() *Pool {
	pool := &Pool{
		NewQueue(),
	}
	return pool
}

// Go works like go func(), but goroutines are pooled for reusing.
// This strategy can avoid runtime.moreStack, because pooled goroutine is already enlarged.
func (pool *Pool) Go(f func()) {
	g := pool.get()
	g.ch <- f
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

func (g *goroutine) workLoop(pool *Pool) {
	for {
		work := <-g.ch
		work()
		pool.Put(g)
	}
}
