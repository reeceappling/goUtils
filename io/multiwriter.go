package io

import (
	"errors"
	"github.com/reeceappling/goUtils/v2/utils"
	"github.com/reeceappling/goUtils/v2/utils/channels"
	"io"
	"math"
	"sync"
)

type MultiWriter interface {
	io.Writer
	Add(io.Writer) error
	Remove(io.Writer) error
}

type flexibleMultiWriter struct {
	writeFunc func(stopEarly bool, toWrite []byte, writeTo []io.Writer) (nMin int, err error)
	writers   utils.Set[io.Writer]
	stopEarly bool
	lock      sync.Mutex
}

func (mw *flexibleMultiWriter) Write(p []byte) (n int, err error) {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	return mw.writeFunc(mw.stopEarly, p, mw.writers.ToSlice())
}

func (mw *flexibleMultiWriter) Add(w io.Writer) error {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	mw.writers.Add(w)
	return nil // TODO: is this ok?
}

func (mw *flexibleMultiWriter) Remove(w io.Writer) error {
	mw.lock.Lock()
	defer mw.lock.Unlock()
	mw.writers.Remove(w)
	return nil // TODO: is this ok?
}

func NewParallelMultiWriter(stopEarly bool, initialWriters ...io.Writer) MultiWriter {
	return &flexibleMultiWriter{
		writeFunc: writeInParallel,
		writers:   utils.SetOf(initialWriters),
		stopEarly: stopEarly,
		lock:      sync.Mutex{},
	}
}

func NewSeriesMultiWriter(stopEarly bool, initialWriters ...io.Writer) MultiWriter {
	return &flexibleMultiWriter{
		writeFunc: writeInSeries,
		writers:   utils.SetOf(initialWriters),
		stopEarly: stopEarly,
		lock:      sync.Mutex{},
	}
}

func writeInSeries(stopEarly bool, toWrite []byte, writeTo []io.Writer) (nMin int, err error) {
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

func writeInParallel(stopEarly bool, toWrite []byte, writeTo []io.Writer) (nMin int, err error) {
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

func writeAsync(p []byte, w io.Writer) <-chan writerOutput {
	resultChan := make(chan writerOutput)
	go func() {
		defer close(resultChan)
		wn, werr := w.Write(p)
		resultChan <- writerOutput{wn, werr}
	}()
	return resultChan
}

type writerOutput struct {
	n   int
	err error
}
