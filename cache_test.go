package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.001
// @date    2018-12-19

import (
	"strconv"
	"testing"
	"time"
)

func TestCache(t *testing.T) {

	cfg := &Config{Nodes: 16, Size: 1600, TTL: 600}

	cache := New(cfg)

	if cache.Size() != 0 {
		t.Fatal("New create not empty cache")
	}

	cache.Flush()

	for i := 0; i < 100; i++ {
		cache.Set(strconv.Itoa(i), i)
	}

	if cache.Size() != 100 {
		t.Fatal("invalid cache size")
	}

	val := cache.Get("50")

	if val == nil || val.(int) != 50 {
		t.Fatal("Set not work")
	}

	fn := func(k string) interface{} {
		return time.Now().Unix()
	}

	res := cache.Fetch("1010101", fn)
	if res == nil {
		t.Fatal("fetch not work")
	}

	prev := res.(int64)
	if prev == 0 {
		t.Fatal("fetch return zero")
	}

	res = cache.Fetch("1010101", fn)
	if res == nil || res.(int64) != prev {
		t.Fatal("fetch not work")
	}

	cache = New(&Config{Nodes: 1, TTL: 3600, Size: 128})

	for i := 0; i < 12; i++ {
		cache.Set(strconv.Itoa(i), i+13)
	}

	for i := 11; i > -1; i-- {

		if cache.Get(strconv.Itoa(i)).(int) != i+13 {
			t.Fatal("cache with one node not worl")
		}
	}
}
