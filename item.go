package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.002
// @date    2019-10-15

import (
	"time"
)

type item struct {
	data   interface{} // value
	expire int64       // expire time
}

func (i *item) isAlive() bool {
	return time.Now().Unix() < i.expire
}
