package nodenum

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-12-06

import (
	"testing"
)

func TestNodeNum(t *testing.T) {
	n := New(0)

	tF := func(key string, wait int) {
		num := n.Get(key)
		if num != wait {
			t.Fatalf("Invalid answer for key: %s wait: %d return: %d", key, wait, num)
		}
	}

	tF("123", 0)
	tF("1234", 0)

	n = New(8)

	tF("123", 2)
	tF("1234", 3)
	tF("12345", 4)
	tF("666", 4)
	tF("777", 4)
	tF("5343", 7)
}
