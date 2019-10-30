package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.006
// @date    2019-10-30

import (
	"sync"
)

const (
	coeff_limit     float64 = 1.1
	coeff_threshold float64 = 0.8
)

type node struct {
	sync.RWMutex
	data      map[string]*item
	prev      map[string]*item
	limit     int
	total     int
	prevTotal int
	threshold int
}

func makeNode(limit int) *node {
	return &node{
		data:      make(map[string]*item, int(float64(limit)*coeff_limit)),
		prev:      map[string]*item{},
		limit:     limit,
		threshold: int(float64(limit) * coeff_threshold),
		total:     0,
	}
}

func (n *node) get(key string) interface{} {
	n.RLock()
	defer n.RUnlock()

	if v, h := n.data[key]; h {

		if v.isAlive() {
			return v.data
		}

		return nil
	}

	if v, h := n.prev[key]; h {

		if v.isAlive() {
			return v.data
		}
	}

	return nil
}

func (n *node) set(key string, value interface{}, before int64) {
	n.Lock()
	defer n.Unlock()

	if old, has := n.data[key]; has {
		old.expire = before
		old.data = value
	} else {
		n.gc()
		n.total++
		n.data[key] = &item{data: value, expire: before}
	}
}

func (n *node) gc() {

	if n.total >= n.threshold {
		n.prev = n.data
		n.data = make(map[string]*item, int(float64(n.limit)*coeff_limit))
		n.prevTotal = n.total
		n.total = 0
	}
}

func (n *node) delete(key string) {
	n.Lock()
	defer n.Unlock()

	if v, has := n.data[key]; has {
		v.expire = 0
	} else if v, has := n.prev[key]; has {
		v.expire = 0
	}
}

func (n *node) size() int {
	n.RLock()
	defer n.RUnlock()

	return n.total + n.prevTotal
}

func (n *node) flush() {
	n.Lock()
	defer n.Unlock()

	n.data = make(map[string]*item, int(float64(n.limit)*coeff_limit))
	n.prev = map[string]*item{}
	n.total = 0
	n.prevTotal = 0
}

func max(a1, a2 int64) int64 {
	if a1 > a2 {
		return a1
	}

	return a2
}

func (n *node) incby(key string, shift int64, before int64) int64 {
	n.Lock()
	defer n.Unlock()

	val := max(shift, 0)

	if old, has := n.data[key]; has {

		if old.isAlive() {

			if v, ok := old.data.(int64); ok {
				val = max(v+shift, 0)
				old.data = val
			} else {
				old.data = val
			}

		} else {
			old.data = val
		}

		old.expire = before

	} else if old, has := n.prev[key]; has {

		if old.isAlive() {

			if v, ok := old.data.(int64); ok {
				val = max(v+shift, 0)
				old.data = val
			} else {
				old.data = val
			}

		} else {
			old.data = val
		}

		old.expire = before

		n.gc()
		n.total++

		n.data[key] = old

	} else {
		n.gc()
		n.total++
		n.data[key] = &item{data: val, expire: before}
	}

	return val
}
