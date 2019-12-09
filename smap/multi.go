package smap

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-12-09

import (
	"github.com/belfinor/lcache/nodenum"
)

type multi struct {
	hash  *nodenum.NodeNum
	nodes []*node
}

func (m *multi) Get(key string) (string, bool) {
	num := m.hash.Get(key)
	return m.nodes[num].Get(key)
}

func (m *multi) Set(key string, val string) {
	num := m.hash.Get(key)
	m.nodes[num].Set(key, val)
}

func (m *multi) Fetch(key string, fn FetchFunc) (string, bool) {
	num := m.hash.Get(key)
	return m.nodes[num].Fetch(key, fn)
}

func (m *multi) Delete(key string) {
	num := m.hash.Get(key)
	m.nodes[num].Delete(key)
}

func (m *multi) Flush() {
	for _, n := range m.nodes {
		n.Flush()
	}
}

func (m *multi) IncBy(key string, inc int64) int64 {
	num := m.hash.Get(key)
	return m.nodes[num].IncBy(key, inc)
}
