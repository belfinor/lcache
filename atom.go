package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2018-07-30

import (
	"sync"
	"time"
)

type Atom struct {
	sync.RWMutex
	ttl    int64
	expire int64
	data   interface{}
}

type ATOM_FETCH_FUNC func() interface{}

func NewAtom(ttl int64) *Atom {
	return &Atom{
		ttl:    ttl,
		expire: 0,
		data:   nil,
	}
}

func (a *Atom) Get() interface{} {
	if a == nil {
		return nil
	}

	a.RLock()
	defer a.RUnlock()

	if a.expire < time.Now().Unix() {
		return nil
	}

	return a.data
}

func (a *Atom) Set(v interface{}) {
	if a == nil {
		return
	}

	a.Lock()
	defer a.Unlock()

	a.expire = time.Now().Unix() + a.ttl
	a.data = v
}

func (a *Atom) Fetch(fn ATOM_FETCH_FUNC) interface{} {
	v := a.Get()
	if v != nil {
		return v
	}

	v = fn()
	if v != nil {
		a.Set(v)
	}

	return v
}

func (a *Atom) Reset() {
	a.Set(nil)
}
