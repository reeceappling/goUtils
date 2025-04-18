package io

import (
	"errors"
	"github.com/reeceappling/goUtils/v2/utils"
	"github.com/reeceappling/goUtils/v2/utils/channels"
	"io"
	"math"
	"sync"
)

type MultiWriteCloser interface {
	io.WriteCloser
	Add(io.WriteCloser) error
	Remove(io.WriteCloser) error
}

type flexibleMultiWriteCloser struct {
	closed    bool
	writeFunc func(stopEarly bool, toWrite []byte, writeTo []io.WriteCloser) (nMin int, err error)
	writers   utils.Set[io.WriteCloser]
	stopEarly bool
	lock      sync.Mutex
}

func (mw *flexibleMultiWriteCloser) Write(p []byte) (n int, err error) {
	if mw.closed {
		return 0, errors.New("cannot write to closed writer")
	}
	mw.lock.Lock()
	defer mw.lock.Unlock()

	return mw.writeFunc(mw.stopEarly, p, mw.writers.ToSlice())
}

func (mw *flexibleMultiWriteCloser) Add(w io.WriteCloser) error {
	if mw.closed {
		return errors.New("cannot add to closed writer")
	}
	mw.lock.Lock()
	defer mw.lock.Unlock()
	mw.writers.Add(w)
	return nil // TODO: is this ok?
}

func (mw *flexibleMultiWriteCloser) Remove(w io.WriteCloser) error {
	if mw.closed {
		return errors.New("cannot remove from closed writer")
	}
	mw.lock.Lock()
	defer mw.lock.Unlock()
	mw.writers.Remove(w)
	return nil // TODO: is this ok?
}

func (mw *flexibleMultiWriteCloser) Close() error {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	if mw.closed {
		return nil
	}
	var err error = nil
	mw.closed = true
	for _, wc := range mw.writers.ToSlice() {
		err = errors.Join(err, wc.Close())
	}
	return err
}

func NewParallelMultiWriteCloser(stopEarly bool, initialWriters ...io.WriteCloser) MultiWriteCloser {
	return &flexibleMultiWriteCloser{
		writeFunc: writeCloserInParallel,
		writers:   utils.SetOf(initialWriters),
		stopEarly: stopEarly,
		lock:      sync.Mutex{},
	}
}

func NewSeriesMultiWriteCloser(stopEarly bool, initialWriters ...io.WriteCloser) MultiWriteCloser {
	return &flexibleMultiWriteCloser{
		writeFunc: writeCloserInSeries,
		writers:   utils.SetOf(initialWriters),
		stopEarly: stopEarly,
		lock:      sync.Mutex{},
	}
}

func writeCloserInSeries(stopEarly bool, toWrite []byte, writeTo []io.WriteCloser) (nMin int, err error) {
	nMin, err = math.MaxInt, nil
	if len(writeTo) == 0 {
		return 0, nil
	}
	for _, w := range writeTo {
		wn, werr := w.Write(toWrite)
		if werr != nil {
			if stopEarly {
				return wn, werr
			}
			err = errors.Join(err, werr)
		}
		nMin = min(nMin, wn)
	}
	return nMin, err
}

func writeCloserInParallel(stopEarly bool, toWrite []byte, writeTo []io.WriteCloser) (nMin int, err error) {
	nMin, err = math.MaxInt, nil
	if len(writeTo) == 0 {
		return 0, nil
	}
	resultChans := make([]<-chan writerOutput, len(writeTo))
	for i, w := range writeTo {
		resultChans[i] = writeAsync(toWrite, w)
	}
	for result := range channels.Multiplex(resultChans) {
		err = errors.Join(err, result.err)
		if err != nil && stopEarly {
			return result.n, result.err // TODO: is this ok?
		}
		nMin = min(nMin, result.n)
	}
	return nMin, err
}

func writeAsyncCloser(p []byte, w io.WriteCloser) <-chan writerOutput { // TODO: ?
	resultChan := make(chan writerOutput)
	go func() {
		defer close(resultChan)
		wn, werr := w.Write(p)
		resultChan <- writerOutput{wn, werr}
	}()
	return resultChan
}
