package ring


// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.000
// @date    2018-01-16


import (
  "fmt"
  "testing"
)


func Test( t *testing.T ) {
  r := New(5)

  if r.Size() != 0 || r.Shift() != "" {
    t.Fatal( "empty cache error" )
  }

  r.Add("1")
  r.Add("2")

  if r.Size() != 2 || r.First() != "1" || r.Get(0) != "1" || r.Get(1) != "2" {
    t.Fatal( "Add error" )
  }

  for i := 0 ; i < 4 ; i++ {
    r.Add( fmt.Sprintf( "%d", i + 3 ) )
  }

  if r.Size() != 5 || r.First() != "1" || r.Get(4) != "5" {
    t.Fatal( "add 4 items error" )
  }

  if r.Shift() != "1" && r.Size() != 4 {
    t.Fatal( "Shift not work" )
  }

  if !r.Add("6") || r.Add("7") || r.Size() != 5 || r.Get(4) != "6" {
    t.Fatal( "Add after Shift error" )
  }
}

