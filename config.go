package lcache


// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.001
// @date    2018-01-16


type Config struct {
  Nodes       int `json:"nodes"`
  Size        int `json:"size"`
  TTL         int `json:"ttl"`
  Clean       int `json:"clean"`
  InputBuffer int `json:"buffer"`
}


func DefaultConfig() *Config {
  return &Config{
    Nodes:       16,
    Size:        10000,
    TTL:         3600,
    Clean:       500,
    InputBuffer: 1000,
  }
}

