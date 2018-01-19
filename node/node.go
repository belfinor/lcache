package node


// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.001
// @date    2018-01-16


import (
  "sync"
)


type Node struct {
  sync.Mutex
  data map[string]*Item
}


func New() *Node {
  return &Node{ data: make( map[string]*Item ) }
}


func (n *Node) Get( key string ) interface{} {
  n.Lock()
  defer n.Unlock()

  if v, h := n.data[key] ; h && v.IsAlive() {
    return v.Data
  }

  return nil
}


func (n *Node) Set( key string, value interface{}, before int64 ) bool {
  n.Lock()
  defer n.Unlock()

  old := n.data[key]
  if old != nil {
    old.Expire = before
    old.Data   = value
  } else {
    n.data[key] = &Item{ Data: value, Expire: before }
  }
  return old == nil
}


func (n *Node) Expire( key string ) {
  n.Lock()
  defer n.Unlock()

  if old, h := n.data[key] ; h {
    old.Expire = -1
  }
}


func (n *Node) Delete( key string ) {
  n.Lock()
  defer n.Unlock()
  delete( n.data, key )
}


func (n *Node) Size() int {
  n.Lock()
  defer n.Unlock()
  return len(n.data)
}


func (n *Node) Flush() {
  n.Lock()
  defer n.Unlock()
  n.data = make( map[string]*Item )
}

