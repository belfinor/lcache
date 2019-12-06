package dmap

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-12-06

import (
	"sync"
)

const (
	coeff_limit     float64 = 1.2
	coeff_threshold float64 = 0.8
)

type node struct {
	sync.RWMutex
	data      map[string]int64
	prev      map[string]int64
	limit     int
	total     int
	prevTotal int
	threshold int
}

func makeNode(limit int) *node {
	return &node{
		data:      make(map[string]int64, int(float64(limit)*coeff_limit)),
		prev:      map[string]int64{},
		limit:     limit,
		threshold: int(float64(limit) * coeff_threshold),
		total:     0,
	}
}

func (n *node) Get(key string) (int64, bool) {
	n.RLock()
	defer n.RUnlock()

	if v, h := n.data[key]; h {
		return v, true
	}

	if v, h := n.prev[key]; h {
		return v, true
	}

	return 0, false
}

func (n *node) Set(key string, value int64) {
	n.Lock()
	defer n.Unlock()

	if _, has := n.data[key]; has {
		n.data[key] = value
	} else {
		n.gc()
		n.total++
		n.data[key] = value
	}
}

func (n *node) gc() {

	if n.total >= n.threshold {
		n.prev = n.data
		n.data = make(map[string]int64, int(float64(n.limit)*coeff_limit))
		n.prevTotal = n.total
		n.total = 0
	}
}

func (n *node) Delete(key string) {
	n.Lock()
	defer n.Unlock()

	if _, has := n.data[key]; has {
		delete(n.data, key)
		n.total--
	}

	if _, has := n.prev[key]; has {
		delete(n.prev, key)
		n.prevTotal--
	}
}

func (n *node) Size() int {
	n.RLock()
	defer n.RUnlock()

	return n.total + n.prevTotal
}

func (n *node) Flush() {
	n.Lock()
	defer n.Unlock()

	n.data = make(map[string]int64, int(float64(n.limit)*coeff_limit))
	n.prev = map[string]int64{}
	n.total = 0
	n.prevTotal = 0
}

func max(a1, a2 int64) int64 {
	if a1 > a2 {
		return a1
	}

	return a2
}

func (n *node) IncBy(key string, shift int64) int64 {
	n.Lock()
	defer n.Unlock()

	val := max(shift, 0)

	if old, has := n.data[key]; has {

		val = max(old+shift, 0)
		n.data[key] = val

	} else if old, has := n.prev[key]; has {

		n.gc()
		n.total++

		val = max(old+shift, 0)
		n.data[key] = val

	} else {
		n.gc()
		n.total++
		n.data[key] = val
	}

	return val
}

func (n *node) Fetch(key string, fn FetchFunc) (int64, bool) {

	if res, has := n.Get(key); has {
		return res, has
	}

	if res, has := fn(key); has {
		n.Set(key, res)
		return res, true
	}

	return 0, false
}
