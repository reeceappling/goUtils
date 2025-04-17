package channels

import "sync"

const (
	SizeBuffer = 1000 // performance testing showed decrease in throughput over 1k buffer
	// TODO: keep or remove next line?
	//busyTimeout = 1 * time.Millisecond // the amount of time to attempt to push data to a worker before giving up
)

// Drain removes all items from all provided channels
func Drain[T any](c ...<-chan T) { // TODO: test?
	go func() {
		for range c {
		}
	}()
}

// Fanout makes one input channel that feeds exact copies into all the input channels.
// Do NOT use this for channels transporting pointers (or slices, etc) if multiple channels plan to modify the pointers
func Fanout[T any](outputs []chan<- T) chan<- T { // TODO: test
	input := make(chan T, SizeBuffer)

	var wg sync.WaitGroup
	wg.Add(len(outputs))
	go func() { // background the wait so the output channel can return
		wg.Wait() // wait for all inputs to close
		close(input)
	}()

	go func() {
		for v := range input {
			for _, output := range outputs {
				output <- v
			}
		}
		for _, output := range outputs {
			close(output)
			wg.Done()
		}
	}()

	return input
}

// Multiplex takes a slice of input channels and returns a single channel with merged output (many inputs, one output)
func Multiplex[T any](inputs []<-chan T) <-chan T { // TODO: test
	output := make(chan T, SizeBuffer)

	var wg sync.WaitGroup
	wg.Add(len(inputs))
	go func() { // background the wait so the output channel can return
		wg.Wait() // wait for all inputs to close
		close(output)
	}()

	for _, input := range inputs {
		in := input // safe copy of pointer to current channel
		go func() { // background/parallelize each individual input's feed to the output
			defer wg.Done()       // when range completes, decrement the waitgroup
			for msg := range in { // extract all
				output <- msg
			}
		}()
	}

	return output
}
