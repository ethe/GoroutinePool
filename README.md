**Deprecated: On Go 1.11.2 goroutine pool has already slower than natiive `go` keyword.**

# Goroutine Pool
The pool implementation of Go goroutine, it is useful to avoid heavy pressure of `runtime.morestack ` when using goroutine.
It could recycle rather than free goroutines after executed. If the recycled goroutines not be stack reduced yet, then there is no need to request more stack.


## Benchmark
```
goos: darwin
goarch: amd64
pkg: eleme/nex/utils/gp
BenchmarkGoPool-4          	 2000000	       753 ns/op	      32 B/op	       1 allocs/op
BenchmarkGo-4              	 5000000	       318 ns/op	       0 B/op	       0 allocs/op
BenchmarkMorestackPool-4   	 1000000	      2551 ns/op	      64 B/op	       3 allocs/op
BenchmarkMoreStack-4       	  300000	      3758 ns/op	      16 B/op	       1 allocs/op
```

## Usage
```go
package main

import (
  "time"

  "github.com/ethe/GoroutinePool"
)

func main() {
  pool := gp.New(20 * time.Second)  // set idle timeout
  pool.Go(func() {}) // same as `go func(){}()`
}

```

## Lock-free Queue
It also contains a lock-free queue (linked list) minimal implementation.

```go
package main

import (  
  "fmt"

  "github.com/ethe/GoroutinePool"
)

func main() {
  queue := NewQueue()
  queue.Put(1)
  fmt.Println(queue.Get())
}
```
