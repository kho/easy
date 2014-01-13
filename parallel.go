package easy

import (
	"runtime"
)

// Parallel stores the amount of parallelism in terms of maximal
// number of workers. Non-positive values means using current
// GOMAXPROCS at each parallel execution.
type Parallel int

// NumWorkers returns the current number of parallel workers specified
// by p. When p <= 0, the return value may change when GOMAXPROCS is
// changed elsewhere.
func (p Parallel) NumWorkers() int {
	i := int(p)
	if i > 0 {
		return i
	}
	return runtime.GOMAXPROCS(0)
}

// Map maps values from source and write results to the returned
// channel. The results may come in arbitrary order w.r.t. their input
// order from source.
func (p Parallel) Map(f func(interface{}) interface{}, source <-chan interface{}) <-chan interface{} {
	sink := make(chan interface{})
	go func() {
		numWorkers := p.NumWorkers()
		// Input for workers.
		buf := make(chan interface{}, numWorkers)
		// Signal for worker completion.
		done := make(chan struct{})
		// Spawn workers.
		for i := 0; i < numWorkers; i++ {
			go func() {
				for v := range buf {
					sink <- f(v)
				}
				done <- struct{}{}
			}()
		}
		// Send tasks.
		for v := range source {
			buf <- v
		}
		close(buf)
		// Wait until all workers have finished.
		numDone := 0
		for numDone < numWorkers {
			<-done
			numDone++
		}
		close(sink)
	}()
	return sink
}

func (p Parallel) MemMap(f func(interface{}) interface{}, source []interface{}) <-chan interface{} {
	sink := make(chan interface{})
	go func() {
		numWorkers := p.NumWorkers()
		// Input for workers.
		buf := make(chan []interface{}, numWorkers)
		// Signal for worker completion.
		done := make(chan struct{})
		// Spawn workers.
		for i := 0; i < numWorkers; i++ {
			go func() {
				for s := range buf {
					for _, v := range s {
						sink <- f(v)
					}
				}
				done <- struct{}{}
			}()
		}
		// Send tasks.
		jobSize := len(source)
		batchSize := jobSize / (2 * numWorkers)
		if batchSize == 0 {
			batchSize = 1
		}
		i := batchSize
		for ; i < jobSize; i += batchSize {
			buf <- source[i-batchSize : i]
		}
		if i-batchSize < jobSize {
			buf <- source[i-batchSize:]
		}
		close(buf)
		// Wait until all workers have finished.
		numDone := 0
		for numDone < numWorkers {
			<-done
			numDone++
		}
		close(sink)
	}()
	return sink
}
