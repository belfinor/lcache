package node


// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2018-01-15


import (
  "time"
)


type Item struct {
  Data   interface{}  // value
  Expire int64        // epoch live before
}


func (i *Item) IsAlive() bool {
  return time.Now().Unix() < i.Expire
}

