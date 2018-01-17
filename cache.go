package lcache


// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.003
// @date    2018-01-17


import (
  "github.com/belfinor/lcache/node"
  "github.com/belfinor/lcache/ring"
  "github.com/belfinor/Helium/hash/consistent"
  "time"
)


type Cache struct {
  nodes    []*node.Node
  buffer   *ring.Ring
  nodesNum int
  input    chan string
  clean    int
  size     int
  hash     *consistent.Hash
  ttl      int64
}


type FETCH_FUNC func ( key string ) interface{}


func New( cfg *Config ) *Cache {

  c := &Cache{}

  c.nodes    = make( []*node.Node, cfg.Nodes )
  c.nodesNum = cfg.Nodes

  for i, _  := range c.nodes {
    c.nodes[i] = node.New()
  }

  c.input = make( chan string, cfg.InputBuffer )

  c.buffer = ring.New( cfg.Size )
  c.size   = cfg.Size

  c.clean = cfg.Clean

  c.ttl = int64( cfg.TTL )

  c.hash = consistent.New( cfg.Nodes )

  go c.worker()

  return c
}


func (c *Cache) worker() {

  for {

    select {
    case key := <- c.input:

      c.buffer.Add( key )

      if c.buffer.Size() == c.size {
        c.gc()
      }
    }

  }
}


func (c *Cache) gc() {

  for i := 0 ; i < c.clean ; i++ {
    key := c.buffer.Shift()
    if  key == "" {
      continue
    }

    n := c.hash.Get( []byte(key) )
    c.nodes[n].Delete( key )
  }

}


func (c *Cache ) Get( key string ) interface{} {
  n := c.hash.Get([]byte(key))
  return c.nodes[n].Get( key )
}


func (c *Cache) Delete( key string ) {
  n := c.hash.Get( []byte( key ) )
  c.nodes[n].Expire( key )
}


func (c *Cache) Set( key string, value interface{} ) {
  n := c.hash.Get( []byte( key ) )
  if c.nodes[n].Set( key, value, c.ttl + time.Now().Unix() ) {
    c.input <- key
  }
}


func (c *Cache) Fetch( key string, f FETCH_FUNC ) interface{} {
  n := c.hash.Get( []byte(key) )
  v := c.nodes[n].Get( key )

  if v == nil {
    v = f(key)
    if v != nil && c.nodes[n].Set( key, v, c.ttl + time.Now().Unix() ) {
      c.input <- key
    }
  }

  return v
}

