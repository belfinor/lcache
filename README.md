# lcache

Embed caching library for Golang

## Summary

* Require Go >= 1.12
* Written on Go and only available for Go
* Embed library
* Store data as interfaces
* Don't create new go-routine
* Thread-safe
* Store data as key/value
* Support TTL for cache values
* Access through cache
* Caching single object
* Does not support versioned values

## Install

```
go get github.com/belfinor/lcache
```

## Usage

### Import

```
import "github.com/belfinor/lcache"
```

### Create cache object

```
var cache *lcache.Cache = lcache.New(&lcache.Config{TTL: 86400, Size: 102400, Nodes: 16})
```

In this example we create cache object. Cache stores data to 16 buckets. 102400 is expeted maximum key/value pairs in cache. 86400 - value TTL (86400 second = 1 day).

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

Delete key from cache

```
cache.Delete("123")
```

#### Remove all data

```
cache.Flush()
```

#### Inc counter

```
cache.Inc("inc")
cache.Inc("inc")
cache.Inc("inc")

fmt.Println( cache.Get("inc").(int64)) // 3
```

### Caching single object

Атомы кэширования *lcache.Atom* - это упрощенная реализация *lcache.Cache*, которая позволяет закэшировать один объект на заданное время, после истечения времени значение станет равным *nil*. Atom - потокобезопасен. Пример работы с атомом ниже:

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
