package dataflow

import (
	"sync"
)

var (
	filterMu sync.RWMutex
	filters  = make(map[string]NewFilter)
)

func Register(name string, filter NewFilter) {
	filterMu.Lock()
	defer filterMu.Unlock()

	if filters == nil {
		panic("dataflow: Register filters is nil")
	}

	if _, dup := filters[name]; dup {
		panic("dataflow: Register called twice for driver " + name)
	}

	filters[name] = filter
}
