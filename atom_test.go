package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2018-07-30

import (
	"testing"
)

func TestAtom(t *testing.T) {
	cache := NewAtom(3600)
	if cache == nil {
		t.Fatal("NewAtom failed")
	}

	v := cache.Get()
	if v != nil {
		t.Fatal("Empty atom return not nul")
	}

	cache.Set("123")
	str := cache.Get().(string)
	if str != "123" {
		t.Fatal("Atom.Set not work")
	}

	str = cache.Fetch(func() interface{} { return "124" }).(string)
	if str != "123" {
		t.Fatal("Atom.Fetch not work")
	}

	cache.Set(nil)
	v = cache.Get()
	if v != nil {
		t.Fatal("Atom set nil not work")
	}
}
