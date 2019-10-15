package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.002
// @date    2019-10-15

import (
	"testing"
	"time"
)

func TestItem(t *testing.T) {

	i1 := item{expire: time.Now().Unix() - 1}

	if i1.isAlive() {
		t.Fatal("isAlive failed")
	}

	i2 := item{expire: time.Now().Unix() + 5}

	if !i2.isAlive() {
		t.Fatal("isAlive failed")
	}
}
