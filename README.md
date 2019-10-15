# lcache

Embedded caching library for Golang

## Summary

* Require Go >= 1.12
* Written on Go
* Embedded library
* Store data as is (key is string, value - interface)
* Don't create new go-routines
* Thread-safe
* Support TTL for cache values
* Access through cache
* Caching single object
* Does not support versioned values
* Easy to use

## Install

```
go get github.com/belfinor/lcache
```

## Usage

### Import

```
import "github.com/belfinor/lcache"
```

### Create cache

```
var cache *lcache.Cache = lcache.New(&lcache.Config{TTL: 86400, Size: 102400, Nodes: 16})
```

In this example we create cache object. The object can be used in a multi-threaded access environment. The cache stores data in 16 buckets and expected values limit is 102400 (in practice, one and a half times more). 86400 - TTL in seconds for all created objects.

You can create as many cache objects as you need in the program (one cache per object class).


### Cache values

Cache key is a string value.

#### Set/Get

Get - get value from cache; Set - save data to cache. If value not found, return value is *nil*.

```
cache.Set("1", "11")
cache.Set("2",int64(22))

fmt.Println( cache.Get("1").(string), cache.Get("2").(int64), cache.Get("3") == nil )
```

#### Access through cache (Fetch)

```
tm := cache.Fetch( "123", func(k string) interface{} {
  return time.Now().Unix()
  }).(int64)
```

Method *Fetch* return value if it's found in cache. Otherwise call func to get value and store result to cache. If func could not get value then return value is *nil*.

#### Delete

Delete key from cache.

```
cache.Delete("123")
```

#### Remove all data

```
cache.Flush()
```

#### Increment

```
cache.Inc("inc")
cache.Inc("inc")
cache.Inc("inc")

fmt.Println( cache.Get("inc").(int64)) // 3
```

*int64* is type of value.

### Caching single object

If you need to cache only one object then you can use *lcache.Atom*. An object can be used in a multi-threaded access environment. Example how to use below:


```
package main

import "fmt"
import "github.com/belfinor/lcache"

func main() {

  atom := lcache.NewAtom(600) // TTL = 600 second

  fmt.Println( atom.Get() ) // <nil>

  atom.Set("12")

  fmt.Println( atom.Get().(string) ) // 12

  atom.Set(nil)
  fmt.Println( atom.Get() ) // <nil>

  res := atom.Fetch( func() interface{} {
    return []int{1,2,3}
  } )

  fmt.Println( res.([]int) ) // [1 2 3]

}

```

# Used in:

* LiveJournal recommender system and stat services
