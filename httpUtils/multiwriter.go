package httpUtils

import (
	"errors"
	"github.com/reeceappling/goUtils/v2/utils"
	"net/http"
	"sync"
	"time"
)

// TODO: MOSTLY UNTESTED

var (
	_ http.ResponseWriter = &HttpMultiResponseWriter{}
	_ http.ResponseWriter = httpMultiResponseChildWriter{}
)

func NewHttpMultiResponseWriter(key string, mainWriter http.ResponseWriter, pool *RequestPool) *HttpMultiResponseWriter {
	return &HttpMultiResponseWriter{
		requestKey:         key,
		oneHundredStatuses: []int{}, // 1xx only
		statusCode:         0,       // 0 means not set
		originalHeaders:    mainWriter.Header().Clone(),
		mainWriter:         mainWriter,
		writers:            []httpMultiResponseChildWriter{},
		pool:               pool,
	}
}

// HttpMultiResponseWriter is an http.ResponseWriter that utilizes a single writer to write to >=1 child writers
// including errChan for use if a write fails to complete
type HttpMultiResponseWriter struct {
	oneHundredStatuses []int // 1xx only
	statusCode         int
	writing            bool // true when writing has begun?
	sync.Mutex              // To block when things are being actively altered
	requestKey         string
	originalHeaders    http.Header
	mainWriter         http.ResponseWriter
	writers            []httpMultiResponseChildWriter
	pool               *RequestPool
}

// httpMultiResponseChildWriter is a wrapper around an http.ResponseWriter,
// including errChan for use if a write fails to complete
type httpMultiResponseChildWriter struct {
	internalWriter http.ResponseWriter
	errChan        chan error
}

// Header meets http.ResponseWriter
func (dupeClient httpMultiResponseChildWriter) Header() http.Header {
	return dupeClient.internalWriter.Header()
}

// WriteHeader meets http.ResponseWriter
func (dupeClient httpMultiResponseChildWriter) WriteHeader(statusCode int) {
	dupeClient.internalWriter.WriteHeader(statusCode)
}

// Write meets http.ResponseWriter
func (dupeClient httpMultiResponseChildWriter) Write(bytes []byte) (int, error) {
	defer close(dupeClient.errChan)
	out, err := dupeClient.internalWriter.Write(bytes)
	if err != nil {
		dupeClient.errChan <- err // TODO: is this ok or can write be called multiple times??
	}
	return out, err
}

// registerDuplicate registers a duplicate request's writer in the HttpMultiResponseWriter,
// also catching it up to already-written status codes.
func (w *HttpMultiResponseWriter) registerDuplicate(client http.ResponseWriter) (<-chan error, error) {
	if w == nil {
		return nil, errors.New("writer is nil") // TODO: ok?
	}
	for { // TODO: max # loops?
		successFullyLocked := w.TryLock()
		if w.writing == true {
			return nil, errors.New("too late, already writing")
		}
		if successFullyLocked {
			break
		}
		time.Sleep(50 * time.Millisecond) // TODO: ok?
	}
	defer w.Unlock()
	// Write status codes to new listener if other writers have already done it
	for _, code := range w.oneHundredStatuses {
		client.WriteHeader(code)
	}
	if w.statusCode != 0 {
		client.WriteHeader(w.statusCode)
	}
	// create channel and add to multiWriter
	ch := make(chan error)
	w.writers = append(w.writers, httpMultiResponseChildWriter{
		internalWriter: client,
		errChan:        ch,
	})
	return ch, nil
}

// Header meets criteria for http.Header
func (w *HttpMultiResponseWriter) Header() http.Header {
	if w == nil {
		return http.Header{}
	}
	// TODO: DO SOMETHING WEIRD IF ALREADY WRITING?
	return w.mainWriter.Header()
}

// finalizeHeaders calculates out header changes from initial to final,
// then makes those changes on all httpMultiResponseChildWriter clients
func (w *HttpMultiResponseWriter) finalizeHeaders() {
	if w == nil {
		return
	}
	// Get differences in headers
	latestHeader := w.mainWriter.Header().Clone()
	// Adding and replacing headers
	headersTried := utils.Set[string]{}
	for headerKey, newVals := range latestHeader {
		headersTried.Add(headerKey)
		replace := false
		originalVals, exists := w.originalHeaders[headerKey]
		if !exists {
			replace = true
		} else {
			// if not all values are the same, replace
			if len(originalVals) != len(newVals) {
				replace = true
			} else {
				for i, newVal := range newVals {
					if originalVals[i] != newVal {
						replace = true
						break
					}
				}
			}

		}
		if replace {
			for _, writer := range w.writers {
				writer.Header()[headerKey] = newVals
			}
		}
	}
	// Delete any headers that no longer exist
	for headerKey := range w.originalHeaders {
		if !headersTried.Contains(headerKey) {
			for _, writer := range w.writers {
				writer.Header().Del(headerKey)
			}
		}
	}
}

// Write meets http.ResponseWriter
func (w *HttpMultiResponseWriter) Write(bytes []byte) (int, error) {
	if w == nil {
		return 0, errors.New("writer is nil")
	}
	w.writing = true
	w.Lock()
	defer w.Unlock() // TODO: ok?

	// Remove from pool as we begin writing
	w.pool.remove(w.requestKey)

	w.finalizeHeaders()

	// Write to internal writers
	for _, writer := range w.writers {
		_, _ = writer.Write(bytes) // Errors handled by each writer respectively
	}

	return w.mainWriter.Write(bytes)
}

// WriteHeader meets http.ResponseWriter
func (w *HttpMultiResponseWriter) WriteHeader(statusCode int) {
	if w == nil || w.writing {
		// Do nothing if already writing
		return
	}
	if statusCode > 200 {
		if w.statusCode != 0 {
			// do nothing since we already have a >2xx
			return
		}
		w.statusCode = statusCode
	} else {
		w.oneHundredStatuses = append(w.oneHundredStatuses, statusCode)
	}

	w.statusCode = statusCode
	w.mainWriter.WriteHeader(statusCode)
	for _, writer := range w.writers {
		writer.WriteHeader(statusCode)
	}
}
