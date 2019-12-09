package smap

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-12-09

import (
	"errors"

	"github.com/belfinor/kvstring"
	"github.com/belfinor/lcache/nodenum"
)

type FetchFunc func(string) (string, bool)

type SMap interface {
	Get(key string) (string, bool)

	Set(key string, val string)

	Fetch(key string, fn FetchFunc) (string, bool)

	Delete(key string)

	IncBy(key string, inc int64) int64

	Flush()
}

func New(dsn string) (SMap, error) {

	kv, err := kvstring.New(dsn)
	if err != nil {
		return nil, errors.New("invalid params")
	}

	size := kv.GetInt("size", 1024)
	if size < 32 {
		size = 32
	}

	num := kv.GetInt("nodes", 0)
	if num < 2 {
		return makeNode(size), nil
	}

	m := &multi{
		hash: nodenum.New(num),
	}

	for i := 0; i < num; i++ {
		m.nodes = append(m.nodes, makeNode(size/num))
	}

	return m, nil
}
