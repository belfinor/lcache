package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.001
// @date    2018-12-13

import (
	"container/list"
	"time"
)

type item struct {
	data    interface{} // value
	expire  int64       // expire time
	element *list.Element
}

func (i *item) isAlive() bool {
	return time.Now().Unix() < i.expire
}
