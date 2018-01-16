package node


// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2018-01-16


import (
  "testing"
  "time"
)


func Test( t *testing.T ) {

  n := New()

  if n.Get( "1" ) != nil {
    t.Fatal( "Empty node Get error" )
  }

  if !n.Set( "1", "123", time.Now().Unix() + 10 ) {
     t.Fatal( "Set new entry error" )
  }

  v := n.Get( "1" )

  if v == nil || v.(string) != "123" {
    t.Fatal( "Get error" )
  }

  n.Expire( "1" )

  if n.Get("1") != nil || n.Size() != 1 {
    t.Fatal( "Expire not work" )
  }

  if n.Set( "1", "345", time.Now().Unix() + 10  ) {
    t.Fatal( "Overrides existing data error" )
  }

  v = n.Get("1")

  if v == nil || v.(string) != "345" {
    t.Fatal( "Get error" )
  }

  n.Set( "2", []int64{ 1, 2, 3 }, time.Now().Unix() + 2 )

  v = n.Get("2")

  if v == nil || len( v.( []int64 ) ) != 3 || n.Size() != 2 {
    t.Fatal( "BAd store array" )
  }

  n.Delete("1")

  if n.Get("1") != nil || n.Get("2" ) == nil || n.Size() != 1 {
    t.Fatal( "Delete not work" )
  }

  <- time.After( time.Second * 3 )

  if n.Get("2") != nil {
    t.Fatal( "auto expire not work" )
  }
  
}

