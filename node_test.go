package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2018-12-13

import (
	"strconv"
	"testing"
	"time"
)

func TestNode(t *testing.T) {
	n := makeNode(10)

	if n == nil {
		t.Fatal("makeNode failed")
	}

	for i := 0; i < 20; i++ {
		n.set(strconv.Itoa(i), i*2, time.Now().Unix()+10)
	}

	if n.size() != 10 {
		t.Fatal("invalid node size")
	}

	for i := 0; i < 10; i++ {

		if n.get(strconv.Itoa(i)) != nil {
			t.Fatal("get return value for deleted key")
		}
	}

	for i := 10; i < 10; i++ {

		v := n.get(strconv.Itoa(i))

		if v == nil || v.(int) != i*2 {
			t.Fatal("get not found exists key")
		}
	}

	n.delete("15")

	if n.get("15") != nil {
		t.Fatal("delete not work")
	}

	n.set("14", 140, time.Now().Unix()+20)

	if n.queue.Back().Value.(string) != "14" {
		t.Fatal("move element after set not work")
	}

	if n.get("14") == nil || n.get("14").(int) != 140 {
		t.Fatal("set new value for exists key not work")
	}

	before := time.Now().Unix() + 20

	for i := int64(1); i < 20; i++ {
		if n.inc("cnt", before) != i {
			t.Fatal("inc not work")
		}
	}

	n.flush()

	if n.size() > 0 {
		t.Fatal("flush not work")
	}
}
