package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.009
// @date    2019-10-30

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

	if cfg == nil {
		cfg = &Config{
			TTL:   3600,
			Nodes: 4,
			Size:  1024,
		}
	}

	if cfg.Nodes < 1 {
		cfg.Nodes = 1
	}

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
	return c.nodes[n].incby(key, 1, c.ttl+time.Now().Unix())
}

func (c *Cache) IncBy(key string, val int64) int64 {
	n := c.hash.Get([]byte(key))
	return c.nodes[n].incby(key, val, c.ttl+time.Now().Unix())
}

func (c *Cache) DecBy(key string, val int64) int64 {
	n := c.hash.Get([]byte(key))
	return c.nodes[n].incby(key, -val, c.ttl+time.Now().Unix())
}

func (c *Cache) Dec(key string) int64 {
	n := c.hash.Get([]byte(key))
	return c.nodes[n].incby(key, -1, c.ttl+time.Now().Unix())
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
