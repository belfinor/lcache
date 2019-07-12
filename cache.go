package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.007
// @date    2019-07-12

import (
	"time"

	"github.com/belfinor/chash"
)

type Cache struct {
	nodes    []*node
	nodesNum int
	size     int
	hash     *chash.Hash
	ttl      int64
}

type FETCH_FUNC func(key string) interface{}

func New(cfg *Config) *Cache {

	c := &Cache{}

	c.nodes = make([]*node, cfg.Nodes)
	c.nodesNum = cfg.Nodes
	c.size = cfg.Size

	for i := 0; i < cfg.Nodes; i++ {
		c.nodes[i] = makeNode(c.size / cfg.Nodes)
	}

	c.ttl = int64(cfg.TTL)

	c.hash = chash.New(cfg.Nodes)

	return c
}

func (c *Cache) Get(key string) interface{} {
	n := c.hash.Get([]byte(key))
	return c.nodes[n].get(key)
}

func (c *Cache) Delete(key string) {
	n := c.hash.Get([]byte(key))
	c.nodes[n].delete(key)
}

func (c *Cache) Set(key string, value interface{}) {
	n := c.hash.Get([]byte(key))
	c.nodes[n].set(key, value, c.ttl+time.Now().Unix())
}

func (c *Cache) Inc(key string) int64 {
	n := c.hash.Get([]byte(key))
	return c.nodes[n].inc(key, c.ttl+time.Now().Unix())
}

func (c *Cache) Fetch(key string, f FETCH_FUNC) interface{} {
	n := c.hash.Get([]byte(key))
	v := c.nodes[n].get(key)

	if v == nil {
		v = f(key)
		if v != nil {
			c.nodes[n].set(key, v, c.ttl+time.Now().Unix())
		}
	}

	return v
}

func (c *Cache) Size() int {
	cnt := 0

	for _, n := range c.nodes {
		cnt += n.size()
	}

	return cnt
}

func (c *Cache) Flush() {
	for _, n := range c.nodes {
		n.flush()
	}
}
