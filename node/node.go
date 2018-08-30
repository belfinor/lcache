package node

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.003
// @date    2018-08-30

import (
	"sync"
)

type Node struct {
	sync.RWMutex
	data map[string]*Item
}

func New() *Node {
	return &Node{data: make(map[string]*Item)}
}

func (n *Node) Get(key string) interface{} {
	n.RLock()
	defer n.RUnlock()

	if v, h := n.data[key]; h && v.IsAlive() {
		return v.Data
	}

	return nil
}

func (n *Node) Set(key string, value interface{}, before int64) bool {
	n.Lock()
	defer n.Unlock()

	old := n.data[key]
	if old != nil {
		old.Expire = before
		old.Data = value
	} else {
		n.data[key] = &Item{Data: value, Expire: before}
	}
	return old == nil
}

func (n *Node) Inc(key string, before int64) (bool, int64) {
	n.Lock()
	defer n.Unlock()

	var val int64

	old := n.data[key]
	if old != nil {

		if !old.IsAlive() {
			val = 1
		} else {
			v, ok := old.Data.(int64)
			if ok {
				val = v + 1
			} else {
				val = 1
			}
		}

		old.Expire = before
		old.Data = val

	} else {
		val = 1
		n.data[key] = &Item{Data: int64(1), Expire: before}
	}
	return old == nil, val
}

func (n *Node) Expire(key string) {
	n.Lock()
	defer n.Unlock()

	if old, h := n.data[key]; h {
		old.Expire = -1
	}
}

func (n *Node) Delete(key string) {
	n.Lock()
	defer n.Unlock()
	delete(n.data, key)
}

func (n *Node) Size() int {
	n.Lock()
	defer n.Unlock()
	return len(n.data)
}

func (n *Node) Flush() {
	n.Lock()
	defer n.Unlock()
	n.data = make(map[string]*Item)
}
