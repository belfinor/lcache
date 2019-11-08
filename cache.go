package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.010
// @date    2019-11-08

import (
	"time"

	"github.com/belfinor/chash"
	"github.com/belfinor/kvstring"
)

type Cache struct {
	nodes    []*node
	nodesNum int
	size     int
	hash     *chash.Hash
	ttl      int64
}

type FetchFunc func(key string) interface{}

// cache := New("size=1024 nodes=4 ttl=3600")
func New(dsn string) *Cache {

	c := &Cache{}

	kv, err := kvstring.New(dsn)

	if err != nil {
		return nil
	}

	nodes := kv.GetInt("nodes", 4)
	ttl := kv.GetInt64("ttl", 3600)
	size := kv.GetInt("size", 1024)

	if nodes < 1 {
		nodes = 1
	}

	if size < 32 {
		size = 32
	}

	c.nodes = make([]*node, nodes)
	c.nodesNum = nodes
	c.size = size

	for i := 0; i < nodes; i++ {
		c.nodes[i] = makeNode(size / nodes)
	}

	c.ttl = ttl

	c.hash = chash.New(nodes)

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

func (c *Cache) Fetch(key string, f FetchFunc) interface{} {
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
