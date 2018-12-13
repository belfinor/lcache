package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.004
// @date    2018-12-13

import (
	"container/list"
	"sync"
)

type node struct {
	sync.RWMutex
	data  map[string]*item
	queue *list.List
	limit int
}

func makeNode(limit int) *node {
	return &node{
		data:  make(map[string]*item, limit+10),
		queue: list.New(),
		limit: limit,
	}
}

func (n *node) get(key string) interface{} {
	n.RLock()
	defer n.RUnlock()

	if v, h := n.data[key]; h && v.isAlive() {
		return v.data
	}

	return nil
}

func (n *node) set(key string, value interface{}, before int64) {
	n.Lock()
	defer n.Unlock()

	if old, has := n.data[key]; has {
		old.expire = before
		old.data = value
		n.queue.MoveToBack(old.element)
	} else {
		e := n.queue.PushBack(key)
		n.data[key] = &item{data: value, expire: before, element: e}
		n.gc()
	}
}

func (n *node) gc() {
	queue := n.queue

	if queue.Len() > n.limit {
		first := queue.Front()
		key := first.Value.(string)
		delete(n.data, key)
		queue.Remove(first)
	}
}

func (n *node) delete(key string) {
	n.Lock()
	defer n.Unlock()

	if v, has := n.data[key]; has {

		delete(n.data, key)

		n.queue.Remove(v.element)
	}
}

func (n *node) size() int {
	n.RLock()
	defer n.RUnlock()

	return n.queue.Len()
}

func (n *node) flush() {
	n.Lock()
	defer n.Unlock()

	n.data = make(map[string]*item, n.limit+10)
	n.queue.Init()
}

func (n *node) inc(key string, before int64) int64 {
	n.Lock()
	defer n.Unlock()

	val := int64(1)

	if old, has := n.data[key]; has {

		if old.isAlive() {

			if v, ok := old.data.(int64); ok {
				val = v + 1
				old.data = val
			} else {
				old.data = int64(1)
			}

		} else {
			old.data = int64(1)
		}

		old.expire = before
		n.queue.MoveToBack(old.element)

	} else {
		e := n.queue.PushBack(key)
		n.data[key] = &item{data: int64(1), expire: before, element: e}
		n.gc()
	}

	return val
}
