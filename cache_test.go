package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.002
// @date    2019-10-30

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

	for i := 0; i < 10; i++ {
		if cache.Inc("cnt") != int64(i+1) {
			t.Fatal("Inc not work")
		}
	}

	for i := 0; i < 3; i++ {
		if cache.Dec("cnt") != int64(9-i) {
			t.Fatal("Dec not work")
		}
	}

	if cache.IncBy("cnt", 5) != 12 {
		t.Fatal("IncBy not work")
	}

	for i := 0; i < 3; i++ {

		res := cache.DecBy("cnt", 10)

		if i == 0 && res != 2 || i > 0 && res != 0 {
			t.Fatal("DecBy not work")
		}

	}
}
