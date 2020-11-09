// +build !leak

package GoroutinePool

import (
	"sync"
	"testing"
	"time"
)

func TestBasicAPI(t *testing.T) {
	gp := New()
	var wg sync.WaitGroup
	wg.Add(1)
	// cover alloc()
	gp.Go(func() { wg.Done() })
	// cover put()
	wg.Wait()
	// cover get()
	gp.Go(func() {})
}

func TestRace(t *testing.T) {
	gp := New()
	var wg sync.WaitGroup
	begin := make(chan struct{})
	wg.Add(500)
	for i := 0; i < 50; i++ {
		go func() {
			<-begin
			for i := 0; i < 10; i++ {
				gp.Go(func() {
					wg.Done()
				})
				time.Sleep(5 * time.Millisecond)
			}
		}()
	}
	close(begin)
	wg.Wait()
}

func BenchmarkGoPool(b *testing.B) {
	gp := New()
	for i := 0; i < b.N/2; i++ {
		gp.Go(func() {})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gp.Go(dummy)
	}
}

func BenchmarkGo(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		go dummy()
	}
}

func dummy() {
}

func BenchmarkMorestackPool(b *testing.B) {
	gp := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		gp.Go(func() {
			moreStack(false)
			wg.Done()
		})
		wg.Wait()
	}
}

func BenchmarkMoreStack(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			moreStack(false)
			wg.Done()
		}()
		wg.Wait()
	}
}

func moreStack(f bool) {
	var stack [8 * 1024]byte
	if f {
		for i := 0; i < len(stack); i++ {
			stack[i] = 'a'
		}
	}
}
