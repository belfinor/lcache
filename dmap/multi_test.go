package dmap

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-12-06

import (
	"strconv"
	"testing"

	"github.com/belfinor/lcache/nodenum"
)

func TestMulti(t *testing.T) {

	n := &multi{
		hash: nodenum.New(8),
	}

	for i := 0; i < 8; i++ {
		n.nodes = append(n.nodes, makeNode(1024/8))
	}

	tS := func(k string, v int64) {

		n.Set(k, v)
		res, has := n.Get(k)

		if !has {
			t.Fatalf("Get not work for: %s", k)
		}

		if res != v {
			t.Fatalf("Invalid return value for: %s", k)
		}
	}

	tG := func(k string, wait int64, found bool) {
		res, has := n.Get(k)
		if has != found {
			t.Fatalf("Get not work for: %s", k)
		}
		if res != wait {
			t.Fatalf("Get return invalid value for: %s", k)
		}
	}

	for i := int64(0); i < 140; i++ {
		tS(strconv.FormatInt(i, 10), i)
	}

	tS("139", 140)
	tS("103", 1)

	tG("120", 120, true)
	tG("102", 102, true)
	tG("96", 96, true)
	tG("1", 1, true)

	n.Delete("102")
	n.Delete("136")
	n.Delete("36")

	for i := 0; i < 50; i++ {
		n.IncBy("cnt", 1)
		tG("cnt", int64(i)+1, true)
	}

	for i := 0; i < 70; i++ {
		n.IncBy("cnt", -1)
		tG("cnt", max(int64(50-i-1), 0), true)
	}

	if n.IncBy("104", 3) != 107 {
		t.Fatal("Invalid IncBy")
	}

	call := 0

	fnOk := func(k string) (int64, bool) {
		call++
		return 1, true
	}

	fnFailed := func(k string) (int64, bool) {
		return 0, false
	}

	res, ok := n.Fetch("1222", fnOk)
	if res != 1 || ok != true || call != 1 {
		t.Fatal("Fetch failed")
	}

	res, ok = n.Fetch("1222", fnOk)
	if res != 1 || ok != true || call != 1 {
		t.Fatal("Fetch failed")
	}

	res, ok = n.Fetch("300", fnFailed)
	if res != 0 || ok {
		t.Fatal("Fetch failed")
	}

	n.Flush()

}
