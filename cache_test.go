package lcache


// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2018-01-16


import (
  "testing"
  "strconv"
//  "time"
)


func Test( t *testing.T ) {

  cfg := DefaultConfig()
  if cfg == nil {
    t.Fatal( "DefaultConfig error" )
  }

  cfg.Size  = 100
  cfg.Clean = 50
  cfg.InputBuffer = 10

  c := New( cfg )

  if c == nil {
    t.Fatal( "cache New failed" )
  }

  v := c.Get( "test" )

  if v != nil {
    t.Fatal( "Empty cache return value" )
  }

  c.Set( "test", "1" )

  v = c.Get( "test" )
  if v == nil || v.(string) != "1" {
    t.Fatal( "Set not work" )
  }

  for i := 0 ; i < 100000 ; i++ {
    str := strconv.Itoa( i )
    c.Set( str, str )
  }

  v = c.Get( "99950" )
  if v == nil || v.(string) != "99950" {
    t.Fatal( "loop set not work" )
  }

  v = c.Get( "95" )
  if v != nil {
    t.Fatal( "error" )
  }

  v = c.Fetch( "_", func( key string ) interface{} {
    return "AAA"
  } )

  if v == nil || v.(string) != "AAA" {
    t.Fatal( "fetch not work" )
  }

  v = c.Get( "_" )
  if v == nil || v.(string) != "AAA" {
    t.Fatal( "Fetch not save data" )
  }
}

