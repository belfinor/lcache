package lcache

// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.002
// @date    2018-12-13

type Config struct {
	Nodes int `json:"nodes"`
	Size  int `json:"size"`
	TTL   int `json:"ttl"`
}

func DefaultConfig() *Config {
	return &Config{
		Nodes: 16,
		Size:  10000,
		TTL:   3600,
	}
}
