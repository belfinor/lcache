package nodenum

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2019-12-06

import (
	"hash/crc32"
)

type NodeNum struct {
	num int
}

func New(nodes int) *NodeNum {
	h := &NodeNum{
		num: nodes,
	}

	return h
}

func (h *NodeNum) Get(key string) int {

	if h.num < 2 {
		return 0
	}

	v := crc32.ChecksumIEEE([]byte(key))

	return int(v % uint32(h.num))
}
