package ring


// @author  Mikhail Kirillov <mikkirillov@yandex.ru>
// @version 1.001
// @date    2018-01-16


type Ring struct {
  first int
  size  int
  alloc int
  data  []string
}


func New( size int ) *Ring {

  if size < 1 {
    panic( "invalid ring size" )
  }

  r := &Ring{
    first: 0,
    size:  0,
    alloc: size,
    data:  make( []string, size ),
  }

  return r
}


func (r *Ring) First() string {
  if r.size > 0 {
    return r.data[r.first]
  }

  return ""
}


func (r *Ring) Size() int {
  return r.size
}


func (r *Ring) Add( item string ) bool {
  if r.size == r.alloc {
    return false
  }

  r.data[ ( r.first + r.size ) % r.alloc ] = item

  r.size++

  return true
}

func (r *Ring) Shift() string {
  if r.size > 0 {
    v := r.data[ r.first ]
    r.first = ( r.first + 1 ) % r.alloc
    r.size--
    return v
  }

  return ""
}


func (r *Ring) Get( n int ) string {
  if n < 0 || n > r.size {
    return ""
  }

  return r.data[ ( r.first + n ) % r.alloc ]
}


func (r *Ring) Flush() {
  r.size = 0
}

