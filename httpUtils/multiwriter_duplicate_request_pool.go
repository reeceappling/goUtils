package httpUtils

import (
	"net/http"
	"sync"
)

// TODO: COMPLETELY UNTESTED

type RequestPool struct {
	pool      map[string]*HttpMultiResponseWriter // Ptr or no?
	poolMutex sync.RWMutex
}

func NewRequestPool() *RequestPool {
	return &RequestPool{
		pool:      map[string]*HttpMultiResponseWriter{},
		poolMutex: sync.RWMutex{},
	}
}

func (pool *RequestPool) add(key string, w *HttpMultiResponseWriter, doLock bool) {
	if pool == nil {
		return
	}
	if doLock {
		pool.poolMutex.Lock()
		defer pool.poolMutex.Unlock()
	}
	pool.pool[key] = w
}

func (pool *RequestPool) remove(key string) {
	if pool == nil {
		return
	}
	pool.poolMutex.Lock()
	defer pool.poolMutex.Unlock()
	delete(pool.pool, key)
}

func (pool *RequestPool) get(key string, doLock bool) (writer *HttpMultiResponseWriter, exists bool) {
	if pool == nil {
		return nil, false
	}
	if doLock {
		pool.poolMutex.RLock()
		defer pool.poolMutex.RUnlock()
	}
	writer, exists = pool.pool[key]
	return writer, exists
}

func (pool *RequestPool) getOrRegister(key string, w http.ResponseWriter) (http.ResponseWriter, <-chan error) {
	if pool == nil {
		return w, nil // TODO: ok?
	}
	pool.poolMutex.Lock()
	defer pool.poolMutex.Unlock()
	writer, exists := pool.get(key, false)
	if !exists {
		return pool.newWriter(key, w, false), nil
	}
	errChan, err := writer.registerDuplicate(w)
	if err != nil {
		return w, nil
	}
	return nil, errChan // TODO: ok, could this fail? Maybe ensure with a mutex
}

func (pool *RequestPool) newWriter(key string, mainWriter http.ResponseWriter, doLock bool) http.ResponseWriter {
	if pool == nil {
		return mainWriter // TODO: ok?
	}
	out := NewHttpMultiResponseWriter(key, mainWriter, pool)
	pool.add(key, out, doLock)
	return out
}

func getPooledWriterOrErrorChannelFunc(deriveKey RequestPoolKeyDeriver) func(http.ResponseWriter, *http.Request, *RequestPool) (http.ResponseWriter, <-chan error) {
	return func(w http.ResponseWriter, r *http.Request, pool *RequestPool) (http.ResponseWriter, <-chan error) {
		key, err := deriveKey(r)
		if err != nil {
			// TODO: WHAT HERE?
		}
		return pool.getOrRegister(key, w)
	}
}

// requestWriterPoolKey is the default RequestPoolKeyDeriver. It does not account for any headers
var requestWriterPoolKey = GeneratePoolKeyDeriver(defaultShortenedUrlKeyDeriver, nilHeaderKeyDeriver, defaultBodyPoolKeyDeriver)

// TODO: TIME OUT THINGS!

// TODO: USE THIS!
func DuplicateRequestPoolMiddleware(pool *RequestPool, handler http.Handler) http.Handler { // TODO: THIS!
	return CustomDuplicateRequestPoolMiddleware(requestWriterPoolKey, pool, handler)
}

// TODO: USE THIS!
func CustomDuplicateRequestPoolMiddleware(deriveKey RequestPoolKeyDeriver, pool *RequestPool, handler http.Handler) http.Handler { // TODO: THIS!
	getWriterOrChannel := getPooledWriterOrErrorChannelFunc(deriveKey)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mainWriter, duplicateListener := getWriterOrChannel(w, r, pool)
		if mainWriter != nil {
			handler.ServeHTTP(mainWriter, r)
			return
		}
		if err := <-duplicateListener; err != nil {
			// TODO: THIS??? LOG????
		}
	})
}
